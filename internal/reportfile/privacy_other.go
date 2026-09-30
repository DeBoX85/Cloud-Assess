//go:build !windows

package reportfile

func preserveReplacementPrivacy(string, string) error { return nil }
