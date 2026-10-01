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

// defaults returns the unchanged built-in identity.
func defaults() Branding {
	return Branding{
		ProductName:      "Cloud Assess",
		ShortName:        "Cloud Assess",
		CLIName:          "cloud-assess",
		ReportTitle:      "Azure Cloud Assessment",
		ReportFilePrefix: "cloud_assessment",
		WebsiteURL:       "https://github.com/DeBoX85/Cloud-Assess",
		SupportURL:       "https://github.com/DeBoX85/Cloud-Assess/issues",
		LogoAltText:      "Cloud Assess logo",
	}
}

// Default returns a copy of the immutable identity resolved at process startup.
// CLI startup checks EmbeddedError before constructing commands or authenticating.
func Default() Branding {
	if embeddedErr != nil {
		panic(embeddedErr)
	}
	return embeddedBrand
}
