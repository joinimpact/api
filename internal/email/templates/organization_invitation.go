package templates

import (
	"fmt"
	"strings"
)

const organizationInvitationTemplate = `
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
              An invitation for you
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
            Hey, {{name}}. You're receiving this email because you've been invited to join the organization <b>{{organizationName}}</b> on Impact.
            If you didn't expect this email, please ignore it and the invite will expire.
          </p>
        </td>
      </tr>
      <tr>
        <td>
          <a
            style="color: initial;"
            href="https://joinimpact.org/dashboard/user/organizations/{{organizationID}}/invites/{{inviteID}}"
            >Click here to join {{organizationName}}</a
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

// OrganizationInvitationTemplate generates and returns a reset password email with the
// provied name, email and key.
func OrganizationInvitationTemplate(name, organizationName string, organizationID, inviteID int64) string {
	template := organizationInvitationTemplate

	// Replace the template variables with the provided params.
	template = strings.Replace(template, `{{name}}`, name, -1)
	template = strings.Replace(template, `{{organizationName}}`, organizationName, -1)
	template = strings.Replace(template, `{{organizationID}}`, fmt.Sprintf("%d", organizationID), -1)
	template = strings.Replace(template, `{{inviteID}}`, fmt.Sprintf("%d", inviteID), -1)

	// Return the HTML string.
	return template
}
