package templates

import "testing"

const expectedResetPasswordOutput = `
<!DOCTYPE html>
<html>
  <body>
    <table
      style="
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
          Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji',
          'Segoe UI Symbol';
        font-size: 16px;
        color: rgb(10, 31, 51);
        width: 100%;
        padding: 30px;
        margin: 15px;
        border-radius: 12px;
        border: 1px solid #e0e7ee;
        max-width: 500px;
        box-shadow: 0px 2px 4px rgba(10, 31, 51, 0.04);
      "
    >
      <tr style="display: block; padding-bottom: 30px;">
        <td>
          <img
            width="156"
            src="https://impact-cdn.sfo2.digitaloceanspaces.com/impact-logo.png"
          />
        </td>
      </tr>
      <tr
        style="
          display: block;
          padding: 8px 0px;
          padding-top: 30px;
          border-top: 1px solid #e0e7ee;
        "
      >
        <td>
          <center>
            <h3 style="font-size: 1.25rem; font-weight: bold; margin: 0;">
              Password Reset
            </h3>
          </center>
        </td>
      </tr>
      <tr style="display: block; padding: 12px 0px; padding-top: 0px;">
        <td>
          <p
            style="
              margin: 0;
              line-height: 1.5;
              color: rgb(69, 94, 117);
              max-width: 360px;
            "
          >
            Hey, Yury. You're receiving this email because you requested to
            reset your password on Impact. If you did not request this, please
            ignore this email and consider changing your password. If you did
            request this, please click the link below to reset your password.
          </p>
        </td>
      </tr>
      <tr>
        <td>
          <a
            style="color: initial;"
            href="https://dev.joinimpact.org/auth/reset/12345678910"
            >Click here to reset your password</a
          >
        </td>
      </tr>
      <tr style="display: block; padding-top: 16px;">
        <td>
          <p
            style="
              margin: 0;
              line-height: 1.5;
              color: rgb(69, 94, 117);
              max-width: 360px;
            "
          >
            Love,
            <br />
            Impact
          </p>
        </td>
      </tr>
    </table>
  </body>
</html>
`

// TestResetPasswordTemplate tests the reset password template.
func TestResetPasswordTemplate(t *testing.T) {
	// Generate a template.
	resetPassword := ResetPasswordTemplate("Yury", "yury@joinimpact.org", "12345678910")

	if resetPassword != expectedResetPasswordOutput {
		// There was a mismatch.
		t.Fatal("template not generated as expected")
	}
}
