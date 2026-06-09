package badge

var declarativeBadges = []*Badge{
	{ID: "works-on-my-machine", Name: map[string]string{"zh": "在我机器上能运行", "en": "Works On My Machine"}, Subtitle: map[string]string{"zh": "态度即正义", "en": "Attitude is everything"}, Icon: "checkmark", Type: Declarative, Rarity: Bronze},
	{ID: "read-not-reply", Name: map[string]string{"zh": "已读不回", "en": "Read Not Reply"}, Subtitle: map[string]string{"zh": "Review了你的PR，然后…没有然后了", "en": "Reviewed your PR, then... nothing"}, Icon: "eye", Type: Declarative, Rarity: Bronze},
	{ID: "stackoverflow-courier", Name: map[string]string{"zh": "Stack Overflow搬运工", "en": "Stack Overflow Courier"}, Subtitle: map[string]string{"zh": "代码从网上来，到网上去", "en": "Code comes from the web, returns to the web"}, Icon: "stack", Type: Declarative, Rarity: Bronze},
	{ID: "todo-collector", Name: map[string]string{"zh": "TODO收藏家", "en": "TODO Collector"}, Subtitle: map[string]string{"zh": "// TODO: fix this later × 50", "en": "// TODO: fix this later × 50"}, Icon: "list", Type: Declarative, Rarity: Bronze},
	{ID: "comment-fundamentalist", Name: map[string]string{"zh": "注释原教旨主义者", "en": "Comment Fundamentalist"}, Subtitle: map[string]string{"zh": "每行代码配三行注释，包括i++", "en": "Three comments per line, including i++"}, Icon: "hash", Type: Declarative, Rarity: Bronze},
	{ID: "copy-paste-engineer", Name: map[string]string{"zh": "复制粘贴工程师", "en": "Copy Paste Engineer"}, Subtitle: map[string]string{"zh": "Ctrl+C / Ctrl+V 是核心技能", "en": "Ctrl+C / Ctrl+V is my core skill"}, Icon: "clipboard", Type: Declarative, Rarity: Bronze},
	{ID: "rubber-duck-master", Name: map[string]string{"zh": "橡皮鸭调试大师", "en": "Rubber Duck Master"}, Subtitle: map[string]string{"zh": "对着鸭子说话就能修bug", "en": "Talk to a duck, fix every bug"}, Icon: "duck", Type: Declarative, Rarity: Silver},
	{ID: "no-friday-deploy", Name: map[string]string{"zh": "周五不部署", "en": "No Friday Deploy"}, Subtitle: map[string]string{"zh": "血的教训换来的铁律", "en": "An iron rule forged in blood"}, Icon: "calendar", Type: Declarative, Rarity: Silver},
	{ID: "force-push-warrior", Name: map[string]string{"zh": "Git Force Push勇士", "en": "Force Push Warrior"}, Subtitle: map[string]string{"zh": "--force 是我的日常", "en": "--force is my daily routine"}, Icon: "zap", Type: Declarative, Rarity: Silver},
	{ID: "meeting-survivor", Name: map[string]string{"zh": "会议幸存者", "en": "Meeting Survivor"}, Subtitle: map[string]string{"zh": "今天开了6个会，写了0行代码", "en": "6 meetings today, 0 lines of code"}, Icon: "users", Type: Declarative, Rarity: Bronze},
}
