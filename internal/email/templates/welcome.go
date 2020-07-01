package templates

import (
	"strings"
)

const welcomeTemplate = `
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
              Welcome to Impact!
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
            Hey, {{name}}. You're receiving this email because you just signed up for an account on Impact.
            Welcome to the club! We're super excited to have you on Impact, and hope that we can provide as much
            value to you as possible. Feel free to reach out to us at <a style="color: initial;" href="mailto:help@joinimpact.org">help@joinimpact.org</a>
            if you need any help.
          </p>
        </td>
      </tr>
      <tr>
        <td>
          <a
            style="color: initial;"
            href="https://joinimpact.org/dashboard"
            >Click here to go to your dashboard</a
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

// WelcomeTemplate generates and returns a welcome email with the
// provied name.
func WelcomeTemplate(name string) string {
	template := welcomeTemplate

	// Replace the template variables with the provided params.
	template = strings.Replace(template, `{{name}}`, name, -1)

	// Return the HTML string.
	return template
}
