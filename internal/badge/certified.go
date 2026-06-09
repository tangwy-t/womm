package badge

var certifiedBadges = []*Badge{
	{ID: "midnight-coder", Name: map[string]string{"zh": "午夜编码者", "en": "Midnight Coder"}, Subtitle: map[string]string{"zh": "月亮不睡我不睡", "en": "The moon doesn't sleep, neither do I"}, Icon: "moon", Type: Certified, Rarity: Rare},
	{ID: "weekend-warrior", Name: map[string]string{"zh": "周末战士", "en": "Weekend Warrior"}, Subtitle: map[string]string{"zh": "工作使我快乐（周末也是）", "en": "Work makes me happy (weekends too)"}, Icon: "sun", Type: Certified, Rarity: Rare},
	{ID: "issue-lord", Name: map[string]string{"zh": "百Issue之主", "en": "Issue Lord"}, Subtitle: map[string]string{"zh": "一切安好……大概", "en": "Everything is fine... probably"}, Icon: "alert", Type: Certified, Rarity: Rare},
	{ID: "docs-master", Name: map[string]string{"zh": "文档仙人", "en": "Docs Master"}, Subtitle: map[string]string{"zh": "代码没写几行，文档写了一本小说", "en": "Few lines of code, a novel of docs"}, Icon: "book", Type: Certified, Rarity: Rare},
	{ID: "pr-bomber", Name: map[string]string{"zh": "PR轰炸机", "en": "PR Bomber"}, Subtitle: map[string]string{"zh": "一天一个PR，医生远离我", "en": "A PR a day keeps the doctor away"}, Icon: "rocket", Type: Certified, Rarity: Rare},
	{ID: "monkey-wrench", Name: map[string]string{"zh": "猴子扳手", "en": "Monkey Wrench"}, Subtitle: map[string]string{"zh": "我来了，CI挂了", "en": "I arrived, CI broke"}, Icon: "wrench", Type: Certified, Rarity: Rare},
	{ID: "archaeologist", Name: map[string]string{"zh": "考古学家", "en": "Archaeologist"}, Subtitle: map[string]string{"zh": "挖出了上古代码", "en": "Unearthed ancient code"}, Icon: "pickaxe", Type: Certified, Rarity: Legendary},
	{ID: "branch-hoarder", Name: map[string]string{"zh": "分支囤积者", "en": "Branch Hoarder"}, Subtitle: map[string]string{"zh": "每个分支都是'马上要合并的'", "en": "Every branch is 'about to be merged'"}, Icon: "git-branch", Type: Certified, Rarity: Rare},
	{ID: "ghost-committer", Name: map[string]string{"zh": "幽灵提交者", "en": "Ghost Committer"}, Subtitle: map[string]string{"zh": "我还活着，只是不想写代码", "en": "I'm alive, just don't want to code"}, Icon: "ghost", Type: Certified, Rarity: Legendary},
	{ID: "polyglot", Name: map[string]string{"zh": "多语言通才", "en": "Polyglot"}, Subtitle: map[string]string{"zh": "什么都会一点，什么都不精", "en": "Jack of all trades, master of none"}, Icon: "globe", Type: Certified, Rarity: Rare},
	{ID: "true-destroyer", Name: map[string]string{"zh": "真·破坏王", "en": "True Destroyer"}, Subtitle: map[string]string{"zh": "连续3次搞挂CI", "en": "Broke CI 3 times in a row"}, Icon: "skull", Type: Certified, Rarity: Legendary},
	{ID: "y2k-hunter", Name: map[string]string{"zh": "千年虫猎人", "en": "Y2K Hunter"}, Subtitle: map[string]string{"zh": "还在跟1999年的代码打交道", "en": "Still dealing with code from 1999"}, Icon: "clock", Type: Certified, Rarity: Legendary},
	{ID: "life-404", Name: map[string]string{"zh": "404人生", "en": "Life 404"}, Subtitle: map[string]string{"zh": "个人主页？不存在的", "en": "Profile page? Not found"}, Icon: "404", Type: Certified, Rarity: Legendary},
	{ID: "commit-anniversary", Name: map[string]string{"zh": "首次提交纪念日", "en": "Commit Anniversary"}, Subtitle: map[string]string{"zh": "写代码这么多年了", "en": "Years of writing code"}, Icon: "trophy", Type: Certified, Rarity: Legendary},
	{ID: "fullstack-victim", Name: map[string]string{"zh": "全栈受害者", "en": "Fullstack Victim"}, Subtitle: map[string]string{"zh": "前端后端运维全都要我干", "en": "Frontend, backend, DevOps — all on me"}, Icon: "layers", Type: Certified, Rarity: Legendary},
}
