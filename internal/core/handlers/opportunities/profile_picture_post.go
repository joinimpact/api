package opportunities

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
	"github.com/oliamb/cutter"
)

// ProfilePicturePost uploads a profile picture to an opportunity.
func ProfilePicturePost(opportunitiesService opportunities.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid file"))
			return
		}
		defer file.Close()

		tmpfile, err := ioutil.TempFile("", "image-upload.*.png")
		if err != nil {
			fmt.Println(err)
			resp.ServerError(w, r, resp.Error(500, "server error"))
			return
		}
		defer os.Remove(tmpfile.Name())

		switch handler.Header.Get("Content-Type") {
		case "image/png", "image/jpeg":
			image, _, err := image.Decode(file)
			if err != nil {
				resp.BadRequest(w, r, resp.Error(400, "invalid file"))
				return
			}

			cropped, err := cutter.Crop(image, cutter.Config{
				Width:   3,
				Height:  1,
				Mode:    cutter.Centered,
				Options: cutter.Ratio,
			})
			if err != nil {
				resp.ServerError(w, r, resp.Error(500, err.Error()))
				return
			}

			err = png.Encode(tmpfile, cropped)
			if err != nil {
				resp.ServerError(w, r, resp.Error(500, "error encoding image"))
				return
			}
			fmt.Println("encoded")
		default:
			resp.BadRequest(w, r, resp.Error(400, "invalid file"))
			return
		}

		f, err := os.Open(tmpfile.Name())
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, "error encoding image"))
			return
		}

		url, err := opportunitiesService.UploadProfilePicture(ctx, opportunityID, f)
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, err.Error()))
			return
		}

		resp.OK(w, r, map[string]interface{}{
			"success":        true,
			"profilePicture": url,
		})
	}
}
