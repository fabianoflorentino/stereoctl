package view

// Profile renders the profile selection screen.
func Profile() string {
	return titleStyle.Render("Profile: DaVinci Resolve (free)") + "\n\n" +
		"→ H.264 video, AAC stereo audio, MP4 container\n\n" +
		dimStyle.Render("[enter] confirm   [esc] back   [q] quit")
}
