package branding

// Branding contains product identity used by presentation and distribution layers.
// Domain schemas and assessment behavior must not depend on these values.
type Branding struct {
	ProductName      string
	ShortName        string
	CLIName          string
	CompanyName      string
	ReportTitle      string
	ReportFilePrefix string
	WebsiteURL       string
	SupportURL       string
	LogoAltText      string
}

// Default returns the working Cloud Assess branding.
func Default() Branding {
	return Branding{
		ProductName:      "Cloud Assess",
		ShortName:        "Cloud Assess",
		CLIName:          "cloud-assess",
		ReportTitle:      "Azure Cloud Assessment",
		ReportFilePrefix: "cloud_assessment",
		LogoAltText:      "Cloud Assess logo",
	}
}
