package render

var icons = map[string]string{
	"checkmark":  `<path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>`,
	"eye":        `<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="1.5" fill="none"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"stack":      `<rect x="4" y="2" width="16" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/><rect x="4" y="10" width="16" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/><rect x="4" y="18" width="16" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"list":       `<line x1="8" y1="6" x2="21" y2="6" stroke="currentColor" stroke-width="1.5"/><line x1="8" y1="12" x2="21" y2="12" stroke="currentColor" stroke-width="1.5"/><line x1="8" y1="18" x2="21" y2="18" stroke="currentColor" stroke-width="1.5"/>`,
	"hash":       `<line x1="4" y1="9" x2="20" y2="9" stroke="currentColor" stroke-width="1.5"/><line x1="4" y1="15" x2="20" y2="15" stroke="currentColor" stroke-width="1.5"/><line x1="10" y1="3" x2="8" y2="21" stroke="currentColor" stroke-width="1.5"/><line x1="16" y1="3" x2="14" y2="21" stroke="currentColor" stroke-width="1.5"/>`,
	"clipboard":  `<rect x="8" y="2" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M16 4h2a2 2 0 012 2v14a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2h2" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"duck":       `<circle cx="12" cy="8" r="5" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M17 10c2 0 4 1 4 3s-2 3-4 3H8" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"calendar":   `<rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="1.5" fill="none"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="1.5"/>`,
	"zap":        `<polygon points="13,2 3,14 12,14 11,22 21,10 12,10" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"users":      `<circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M2 21v-2a4 4 0 014-4h6a4 4 0 014 4v2" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"moon":       `<path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"sun":        `<circle cx="12" cy="12" r="5" stroke="currentColor" stroke-width="1.5" fill="none"/><line x1="12" y1="1" x2="12" y2="3" stroke="currentColor" stroke-width="1.5"/><line x1="12" y1="21" x2="12" y2="23" stroke="currentColor" stroke-width="1.5"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64" stroke="currentColor" stroke-width="1.5"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78" stroke="currentColor" stroke-width="1.5"/><line x1="1" y1="12" x2="3" y2="12" stroke="currentColor" stroke-width="1.5"/><line x1="21" y1="12" x2="23" y2="12" stroke="currentColor" stroke-width="1.5"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36" stroke="currentColor" stroke-width="1.5"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22" stroke="currentColor" stroke-width="1.5"/>`,
	"alert":      `<path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"book":       `<path d="M4 19.5A2.5 2.5 0 016.5 17H20" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"rocket":     `<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 00-2.91-.09z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"wrench":     `<path d="M14.7 6.3a1 1 0 000 1.4l1.6 1.6a1 1 0 001.4 0l3.77-3.77a6 6 0 01-7.94 7.94l-6.91 6.91a2.12 2.12 0 01-3-3l6.91-6.91a6 6 0 017.94-7.94l-3.76 3.76z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"pickaxe":    `<path d="M14.5 2.5L2 15l3 3L17.5 6.5" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M14.5 2.5l7 7-3 3-7-7" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"git-branch": `<line x1="6" y1="3" x2="6" y2="15" stroke="currentColor" stroke-width="1.5"/><circle cx="18" cy="6" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/><circle cx="6" cy="18" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M18 9a9 9 0 01-9 9" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"ghost":      `<path d="M12 2a8 8 0 00-8 8v12l3-3 3 3 3-3 3 3 3-3 3 3V10a8 8 0 00-8-8z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"globe":      `<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" fill="none"/><line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="1"/>`,
	"skull":      `<circle cx="12" cy="10" r="8" stroke="currentColor" stroke-width="1.5" fill="none"/><circle cx="9" cy="9" r="2" fill="currentColor"/><circle cx="15" cy="9" r="2" fill="currentColor"/>`,
	"404":        `<text x="2" y="18" font-family="monospace" font-size="14" font-weight="bold" fill="currentColor">404</text>`,
	"trophy":     `<path d="M8 21h8M12 17v4M7 4h10v4a5 5 0 01-10 0V4z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"layers":     `<polygon points="12,2 2,7 12,12 22,7" stroke="currentColor" stroke-width="1.5" fill="none"/><polyline points="2,17 12,22 22,17" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"clock":      `<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" fill="none"/><polyline points="12,6 12,12 16,14" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
}

func GetIcon(name string) string {
	if icon, ok := icons[name]; ok {
		return icon
	}
	return icons["checkmark"]
}
