        // ========== 数据存储 ==========
        const STORAGE_KEY = 'blog_articles';
        const PROFILE_KEY = 'blog_profile';
        const LINKS_KEY = 'blog_links';
        const PAGES_KEY = 'blog_static_pages';
        const MEDIA_KEY = 'blog_media';
        const SETTINGS_KEY = 'blog_settings';

        // 默认用户信息
        const defaultProfile = {
            username: 'Fu Chengyan',
            desc: '日常落灰的个人博客，擅长面向搜索引擎编程。分享 Golang 开发、AI 和 NAS 折腾经验',
            avatar: null,
            social: [
                { icon: '●●', name: 'WeChat', url: '#wechat' },
                { icon: '⌘', name: 'GitHub', url: 'https://github.com/fcy222fcy?tab=repositories' },
                { icon: '▦', name: 'Tools', url: '#tools' },
                { icon: '◔', name: 'RSS', url: '#rss' }
            ]
        };

        // 初始示例数据
        const defaultArticles = [
            {
                id: 1,
                title: '我的博客主题已开源，欢迎使用',
                summary: '基于 Hugo Theme Stack 打造的开箱即用博客模板，包含多项美化和功能增强，欢迎star和反馈。',
                content: '## 引言\n\n基于 Hugo Theme Stack 打造的开箱即用博客模板，包含多项美化和功能增强。如果你喜欢我的博客样式，不妨试试它！\n\n## 功能特性\n\n- **响应式设计**：完美适配桌面端和移动端\n- **暗色模式**：支持亮色/暗色主题自动切换\n- **多语言支持**：内置中英文切换功能\n- **代码高亮**：支持多种编程语言的代码块高亮显示\n- **图片懒加载**：优化页面加载速度\n\n## 快速开始\n\n```bash\n# 克隆仓库\ngit clone https://github.com/fcy222fcy/hugo-stack-starter.git\n\n# 进入目录\ncd hugo-stack-starter\n\n# 启动本地服务\nhugo server\n```\n\n## 自定义配置\n\n在 `config.yaml` 中修改你的博客信息：\n\n```yaml\ntitle: 你的博客标题\nparams:\n  description: 博客描述\n  author: 作者名\n```\n\n## 预览效果\n\n访问 [在线演示](https://blog.fcyan.xyz) 查看实际效果。\n\n## 许可证\n\n本项目采用 MIT 许可证开源，欢迎自由使用和修改。',
                category: '搭建网站',
                tags: ['Hugo', '开源', '博客'],
                status: 'published',
                views: 252,
                cover: null,
                createdAt: '2026-04-24',
                updatedAt: '2026-04-24'
            },
            {
                id: 2,
                title: 'DesktopSnap：开源轻量、一键恢复，解决 Windows 桌面图标错乱的烦恼',
                summary: '基于 WinUI3 开发的轻量级桌面工具，支持多显示器布局一键备份与恢复，已上架微软商店。',
                content: '## 项目介绍\n\n每次系统更新或者分辨率变化后，桌面图标就会乱成一团？DesktopSnap 可以帮你一键备份和恢复桌面布局！\n\n这是一款基于 WinUI3 开发的轻量级桌面工具，支持多显示器布局一键备份与恢复，已上架微软商店。\n\n## 核心功能\n\n- **一键备份**：快速保存当前桌面图标布局\n- **一键恢复**：随时恢复到之前保存的布局\n- **多显示器支持**：完美兼容多屏幕环境\n- **开机自启**：可设置系统启动时自动运行\n- **轻量无广告**：软件体积小，无任何广告弹窗\n\n## 下载安装\n\n前往 Microsoft Store 搜索 "DesktopSnap" 即可免费下载。\n\n## 技术实现\n\n使用 WinUI3 框架开发，通过 Windows API 获取和设置桌面图标位置。项目完全开源，欢迎贡献代码。\n\n```csharp\n// 获取桌面图标位置示例\nvar icons = DesktopIconManager.GetIconPositions();\nDesktopIconManager.SaveLayout(icons);\n```\n\n## 反馈与贡献\n\n如果遇到问题或有好的建议，欢迎在 GitHub 上提 Issue 或 PR。',
                category: '软件开发',
                tags: ['Windows', 'WinUI3', '开源'],
                status: 'published',
                views: 327,
                cover: null,
                createdAt: '2026-03-12',
                updatedAt: '2026-03-12'
            },
            {
                id: 3,
                title: '用 Cloudflare Workers 免费给博客增加一个带有知识库的 AI 助手',
                summary: '利用 Cloudflare Workers 和 AI 能力，为博客添加智能问答功能，零成本实现。',
                content: '## 前言\n\n想给自己的博客添加一个智能问答助手，但又不想花太多钱？Cloudflare Workers 提供了免费额度，配合 AI 能力，可以零成本实现这个功能。\n\n## 实现原理\n\n1. 使用 Cloudflare Workers 作为后端服务\n2. 接入 AI API 实现智能问答\n3. 利用 Cloudflare KV 存储知识库数据\n4. 前端通过 JavaScript 调用 Worker 接口\n\n## 代码示例\n\n```javascript\n// Worker 入口函数\nexport default {\n  async fetch(request, env) {\n    const url = new URL(request.url);\n    if (url.pathname === "/api/chat") {\n      return handleChat(request, env);\n    }\n    return new Response("Not Found", { status: 404 });\n  }\n}\n\nasync function handleChat(request, env) {\n  const { message } = await request.json();\n  const knowledge = await env.KNOWLEDGE.get("faq", "json");\n  const response = await fetch("https://api.openai.com/v1/chat/completions", {\n    method: "POST",\n    headers: {\n      "Authorization": "Bearer " + env.OPENAI_API_KEY,\n      "Content-Type": "application/json"\n    },\n    body: JSON.stringify({\n      model: "gpt-3.5-turbo",\n      messages: [\n        { role: "system", content: "你是一个博客助手，基于以下知识库回答问题：" + knowledge },\n        { role: "user", content: message }\n      ]\n    })\n  });\n  return Response.json(await response.json());\n}\n```\n\n## 部署步骤\n\n1. 注册 Cloudflare 账号\n2. 创建 Workers 项目\n3. 部署上述代码\n4. 配置自定义域名\n5. 在博客中嵌入对话组件\n\n## 成本分析\n\n- Cloudflare Workers：免费额度足够个人博客使用\n- AI API：根据调用量付费，个人博客场景成本很低\n- KV 存储：免费额度包含 10GB 存储空间\n\n## 总结\n\n通过 Cloudflare Workers，可以低成本地为博客添加智能问答功能，提升用户体验。',
                category: '搭建网站',
                tags: ['Cloudflare', 'AI', 'Workers'],
                status: 'published',
                views: 160,
                cover: null,
                createdAt: '2026-04-21',
                updatedAt: '2026-04-21'
            },
            {
                id: 4,
                title: '清明假期：收拾家务、整理心情',
                summary: '清明假期三天，没有远行，选择在家收拾家务、整理心情，享受难得的宁静时光。',
                content: '## 前言\n\n清明假期三天，没有选择远行，而是在家好好收拾了一下。平时工作忙碌，家里积攒了不少需要整理的东西，趁着假期一次性处理掉。\n\n## 第一天：断舍离\n\n早上睡到自然醒，泡了一杯咖啡，开始翻箱倒柜。把那些一年以上没用过的东西统统清理出来，该扔的扔，该捐的捐。\n\n- 清理了三个纸箱的旧书\n- 整理了衣柜，捐掉了十几件不穿的衣服\n- 清理了电脑里的无用文件\n\n## 第二天：深度清洁\n\n给家里来了一次大扫除：\n- 擦窗户、拖地板\n- 清洗厨房油烟机\n- 整理阳台的杂物\n\n忙了一整天，虽然累但很有成就感。\n\n## 第三天：静心阅读\n\n最后一天，哪儿也没去，在阳台上晒太阳、看书。读完了之前一直想读的《置身事内》，对中国经济有了更深的理解。\n\n## 感悟\n\n有时候，最好的休息不是出去玩，而是让生活回归简单。收拾家务的过程也是整理心情的过程，扔掉不需要的东西，内心也会变得轻松。\n\n假期结束，元气满满地迎接新一周的工作！',
                category: '生活记录',
                tags: ['生活', '随笔'],
                status: 'published',
                views: 100,
                cover: null,
                createdAt: '2026-04-04',
                updatedAt: '2026-04-04'
            },
            {
                id: 5,
                title: '逐渐难以逃离对于 AI 的焦虑',
                summary: 'AI 技术的快速发展让很多人感到焦虑，如何在这波浪潮中找到自己的位置？',
                content: '## 引言\n\n最近和朋友聊天，发现大家都在讨论 AI。有人在学 Prompt Engineering，有人在用 AI 写代码，还有人在研究如何用 AI 做副业。似乎不会用 AI 就要被时代抛弃了。\n\n## 焦虑的来源\n\n1. **技术迭代太快**：ChatGPT、Claude、Gemini，每个月都有新模型发布\n2. **应用场景太多**：写作、编程、设计、视频，AI 无处不在\n3. **替代焦虑**：担心自己的工作会被 AI 取代\n4. **信息过载**：每天都有大量 AI 相关的新闻和教程\n\n## 我的思考\n\n### 焦虑是正常的\n\n面对新技术的冲击，感到焦虑是很正常的。但焦虑不能解决问题，行动才能。\n\n### 找到自己的节奏\n\n不需要追每一个热点，但要保持学习。选择一两个方向深入，比什么都浅尝辄止要好。\n\n### AI 是工具，不是对手\n\n与其担心被 AI 替代，不如思考如何用 AI 提升自己的效率。AI 是工具，关键在于使用它的人。\n\n## 我的实践\n\n- 日常使用 Claude 辅助编程\n- 用 AI 帮忙写技术文档\n- 尝试用 AI 生成博客配图\n\n## 结语\n\n与其焦虑，不如行动。在这个快速变化的时代，保持学习的心态，找到适合自己的节奏，才是最重要的。\n\n> 种一棵树最好的时间是十年前，其次是现在。',
                category: '生活记录',
                tags: ['AI', '焦虑', '思考'],
                status: 'published',
                views: 114,
                cover: null,
                createdAt: '2026-03-16',
                updatedAt: '2026-03-16'
            },
            {
                id: 6,
                title: '开发了一个2026有奖发票抽奖助手，方便合并发票凑齐100元',
                summary: '实用工具开发分享，帮助用户管理发票、合并金额，方便参与有奖发票抽奖活动。',
                content: '## 项目背景\n\n2026年的有奖发票活动又开始了，规则是发票金额满100元可以参与抽奖。但平时开发票都是零散的，要凑齐100元很麻烦。\n\n于是写了一个小工具，帮助管理发票、自动计算合并方案。\n\n## 功能介绍\n\n- **发票录入**：支持手动输入和拍照识别\n- **金额统计**：实时显示总金额和距离100元的差额\n- **智能合并**：自动计算最优合并方案\n- **历史记录**：保存所有发票记录，方便查询\n\n## 技术实现\n\n```go\ntype Invoice struct {\n    ID        int       "json:id"\n    Amount    float64   "json:amount"\n    Date      time.Time "json:date"\n    Seller    string    "json:seller"\n    Merged    bool      "json:merged"\n}\n\n// 合并计算\nfunc calculateMerge(invoices []Invoice, target float64) [][]Invoice {\n    // 使用动态规划计算最优合并方案\n    // ...\n}\n```\n\n## 使用效果\n\n测试了一个月的发票，工具成功帮我合并出了5组100元以上的发票，省去了手动计算的麻烦。\n\n## 后续计划\n\n- 支持 OCR 自动识别发票信息\n- 添加发票有效期提醒\n- 支持导出合并报告\n\n项目已开源，欢迎试用和反馈！',
                category: '软件开发',
                tags: ['工具', '发票', 'Go'],
                status: 'draft',
                views: 97,
                cover: null,
                createdAt: '2026-03-15',
                updatedAt: '2026-03-15'
            }
        ];

        // 默认友情链接
        const defaultLinks = [
            { id: 1, name: 'Hugo', url: 'https://gohugo.io', desc: '世界上最快的网站构建框架' },
            { id: 2, name: 'Go', url: 'https://go.dev', desc: '构建快速、可靠、高效的大规模软件' },
            { id: 3, name: 'GitHub', url: 'https://github.com', desc: '全球开发者构建软件的地方' }
        ];

        // 默认静态页面
        const defaultPages = [
            { id: 1, title: '关于我', slug: 'about', status: 'published', updatedAt: '2026-04-01' },
            { id: 2, title: '友情链接', slug: 'links', status: 'published', updatedAt: '2026-03-20' },
            { id: 3, title: '留言板', slug: 'guestbook', status: 'draft', updatedAt: '2026-03-15' }
        ];

        // 默认系统设置
        const defaultSettings = {
            siteName: 'Fu Chengyan\'s Blog',
            siteDesc: '日常落灰的个人博客，分享 Golang、AI 和 NAS 折腾经验',
            siteUrl: 'https://fcyan.xyz',
            pageSize: 10,
            seoTitle: '',
            seoDesc: '',
            seoKeywords: ''
        };

        // ========== 状态管理 ==========
        let state = {
            articles: [],
            links: [],
            pages: [],
            media: [],
            settings: {},
            currentPage: 'articles',
            editingArticle: null,
            tags: [],
            deleteArticleId: null,
            modalTags: [],
            // 每日一问相关状态
            dailyQuestions: [],
            dailyQuestionComments: [],
            dailyQuestionStats: {
                total_questions: 0,
                enabled_questions: 0,
                total_views: 0,
                total_likes: 0,
                total_comments: 0
            },
            profile: null
        };

        // ========== 主页每日一问小卡片 ==========
        // 模拟数据
        const dqHomeData = {
            id: 1,
            question: '什么是最好的编程语言？',
            answer: '没有最好的编程语言，只有最适合的。Python 适合数据科学，JavaScript 适合 Web 开发，Go 适合后端服务。选择语言要根据具体需求和项目类型。',
            like_count: 42,
            comment_count: 8,
            display_date: '2026-06-15'
        };

        let dqHomeAnswerVisible = false;
        let dqHomeLiked = localStorage.getItem('dq_home_liked') === 'true';

        // 初始化每日一问小卡片
        function initDqHomeCard() {
            const dqHomeDate = document.getElementById('dqHomeDate');
            if (!dqHomeDate) return; // 仪表盘页面无此元素

            // 设置日期
            const date = new Date(dqHomeData.display_date);
            dqHomeDate.textContent = `${date.getMonth() + 1}月${date.getDate()}日`;

            // 设置点赞数和评论数
            const likeCount = document.getElementById('dqHomeLikeCount');
            const commentCount = document.getElementById('dqHomeCommentCount');
            if (likeCount) likeCount.textContent = dqHomeData.like_count;
            if (commentCount) commentCount.textContent = dqHomeData.comment_count;

            // 更新点赞按钮状态
            updateDqHomeLikeButton();
        }

        // 切换答案显示
        function toggleHomeAnswer() {
            const answer = document.getElementById('dqHomeAnswer');
            const btn = document.getElementById('dqHomeToggleBtn');

            if (dqHomeAnswerVisible) {
                answer.style.display = 'none';
                btn.textContent = '查看答案';
                dqHomeAnswerVisible = false;
            } else {
                answer.style.display = 'block';
                btn.textContent = '收起';
                dqHomeAnswerVisible = true;
            }
        }

        // 点赞
        function likeHomeQuestion() {
            if (dqHomeLiked) {
                showToast('已经点过赞了', 'info');
                return;
            }

            dqHomeData.like_count++;
            const likeCountEl = document.getElementById('dqHomeLikeCount');
            if (likeCountEl) likeCountEl.textContent = dqHomeData.like_count;
            localStorage.setItem('dq_home_liked', 'true');
            dqHomeLiked = true;
            updateDqHomeLikeButton();
            showToast('点赞成功', 'success');
        }

        // 更新点赞按钮状态
        function updateDqHomeLikeButton() {
            const btn = document.getElementById('dqHomeLikeBtn');
            if (!btn) return;
            if (dqHomeLiked) {
                btn.classList.add('liked');
            } else {
                btn.classList.remove('liked');
            }
        }

        // 显示评论（示例）
        function showHomeComments() {
            showToast('评论功能', 'info');
        }

        // ========== 初始化 ==========
        function init() {
            const saved = localStorage.getItem(STORAGE_KEY);
            if (saved) {
                state.articles = JSON.parse(saved);
            } else {
                state.articles = defaultArticles;
                saveToStorage();
            }

            // 加载用户信息
            const savedProfile = localStorage.getItem(PROFILE_KEY);
            if (savedProfile) {
                state.profile = JSON.parse(savedProfile);
            } else {
                state.profile = defaultProfile;
                saveProfile();
            }

            // 加载友情链接
            const savedLinks = localStorage.getItem(LINKS_KEY);
            if (savedLinks) {
                state.links = JSON.parse(savedLinks);
            } else {
                state.links = defaultLinks;
                localStorage.setItem(LINKS_KEY, JSON.stringify(state.links));
            }

            // 加载静态页面
            const savedPages = localStorage.getItem(PAGES_KEY);
            if (savedPages) {
                state.pages = JSON.parse(savedPages);
            } else {
                state.pages = defaultPages;
                localStorage.setItem(PAGES_KEY, JSON.stringify(state.pages));
            }

            // 加载媒体库
            const savedMedia = localStorage.getItem(MEDIA_KEY);
            if (savedMedia) {
                state.media = JSON.parse(savedMedia);
            } else {
                state.media = [];
            }

            // 加载系统设置
            const savedSettings = localStorage.getItem(SETTINGS_KEY);
            if (savedSettings) {
                state.settings = JSON.parse(savedSettings);
            } else {
                state.settings = defaultSettings;

            // 加载每日一问数据
            loadDqData();
                localStorage.setItem(SETTINGS_KEY, JSON.stringify(state.settings));
            }

            renderProfile();
            renderAboutPage();
            renderArticleTable();
            renderLinkTable();
            renderPageTable();
            renderMediaGrid();
            loadSettings();
            updateStats();
            updateLinkCount();
            bindEvents();
            initTheme();
            updatePreview();
            initDqHomeCard(); // 初始化每日一问小卡片
        }

        // ========== 主题切换 ==========
        function initTheme() {
            const savedTheme = localStorage.getItem('theme') || 'light';
            document.documentElement.setAttribute('data-scheme', savedTheme);
        }

        function toggleTheme() {
            const current = document.documentElement.getAttribute('data-scheme');
            const next = current === 'light' ? 'dark' : 'light';
            document.documentElement.setAttribute('data-scheme', next);
            localStorage.setItem('theme', next);
        }

        function toggleLanguage() {
            const currentLang = localStorage.getItem('language') || 'zh';
            const nextLang = currentLang === 'zh' ? 'en' : 'zh';
            localStorage.setItem('language', nextLang);
            // 更新按钮文本
            const langText = document.querySelector('#languageToggle span:last-child');
            if (langText) {
                langText.textContent = nextLang === 'zh' ? 'English' : '中文';
            }
            showToast(nextLang === 'zh' ? '已切换到中文' : 'Switched to English', 'success');
        }

        // ========== 保存到 localStorage ==========
        function saveToStorage() {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(state.articles));
        }

        // ========== 保存用户信息 ==========
        function saveProfile() {
            localStorage.setItem(PROFILE_KEY, JSON.stringify(state.profile));
        }

        // ========== 渲染用户信息 ==========
        function renderProfile() {
            const avatarText = document.getElementById('avatarText');
            const username = document.getElementById('sidebarUsername');
            const desc = document.getElementById('sidebarDesc');
            const avatar = document.getElementById('sidebarAvatar');
            const socialContainer = document.getElementById('sidebarSocial');

            if (state.profile.avatar) {
                avatar.innerHTML = `<img src="${state.profile.avatar}" alt="头像"><div class="sidebar-avatar-edit">更换头像</div>`;
            } else {
                avatarText.textContent = state.profile.username.charAt(0).toUpperCase();
            }

            username.textContent = state.profile.username;
            desc.textContent = state.profile.desc;

            // 渲染社交链接
            if (state.profile.social && state.profile.social.length > 0) {
                socialContainer.innerHTML = state.profile.social.map(item =>
                    `<a href="${escapeHtml(item.url)}" class="sidebar-social-item" title="${escapeHtml(item.name)}" target="_blank">${escapeHtml(item.icon)}</a>`
                ).join('');
            }

            // 渲染关于我页面
            renderAboutPage();
        }

        // ========== 渲染关于我页面 ==========
        function renderAboutPage() {
            const aboutAvatarText = document.getElementById('aboutAvatarText');
            const aboutUsername = document.getElementById('aboutUsername');
            const aboutDesc = document.getElementById('aboutDesc');
            const aboutAvatar = document.getElementById('aboutAvatarText')?.parentElement;
            const aboutSocial = document.getElementById('aboutSocial');

            if (aboutUsername) aboutUsername.textContent = state.profile.username;
            if (aboutDesc) aboutDesc.textContent = state.profile.desc;

            if (aboutAvatar) {
                if (state.profile.avatar) {
                    aboutAvatar.innerHTML = `<img src="${state.profile.avatar}" alt="头像" style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;">`;
                } else {
                    aboutAvatarText.textContent = state.profile.username.charAt(0).toUpperCase();
                }
            }

            // 渲染关于我页面的社交链接
            if (aboutSocial && state.profile.social && state.profile.social.length > 0) {
                aboutSocial.innerHTML = state.profile.social.map(item =>
                    `<a href="${escapeHtml(item.url)}" class="sidebar-social-item" title="${escapeHtml(item.name)}" target="_blank">${escapeHtml(item.icon)}</a>`
                ).join('');
            }
        }

        // ========== 渲染编辑弹窗中的社交链接 ==========
        function renderSocialList() {
            const container = document.getElementById('socialList');
            if (!state.profile.social) state.profile.social = [];

            container.innerHTML = state.profile.social.map((item, index) => `
                <div class="social-item" data-index="${index}">
                    <div class="social-item-icon">${escapeHtml(item.icon)}</div>
                    <div class="social-item-fields">
                        <input type="text" class="form-input social-icon-input" value="${escapeHtml(item.icon)}" placeholder="图标" data-index="${index}">
                        <input type="text" class="form-input social-url-input" value="${escapeHtml(item.url)}" placeholder="链接地址" data-index="${index}">
                    </div>
                    <button class="btn btn-danger btn-sm remove-social-btn" data-index="${index}">×</button>
                </div>
            `).join('');
        }

        // ========== 简单 Markdown 转 HTML ==========
        function markdownToHtml(md) {
            if (!md) return '';

            // 保护代码块和行内代码
            const codeBlocks = [];
            const inlineCodes = [];

            let html = md.replace(/```(\w+)?\n([\s\S]*?)```/g, (match, lang, code) => {
                const index = codeBlocks.length;
                codeBlocks.push(`<pre><code class="language-${lang || ''}">${escapeHtml(code)}</code></pre>`);
                return `%%CODEBLOCK_${index}%%`;
            });

            html = html.replace(/`([^`]+)`/g, (match, code) => {
                const index = inlineCodes.length;
                inlineCodes.push(`<code>${escapeHtml(code)}</code>`);
                return `%%INLINECODE_${index}%%`;
            });

            // 处理表格
            html = html.replace(/^(\|.+\|)\n(\|[-:| ]+\|)\n((\|.+\|\n?)*)/gm, (match, header, separator, body) => {
                const headers = header.split('|').filter(c => c.trim()).map(c => `<th>${c.trim()}</th>`).join('');
                const rows = body.trim().split('\n').map(row => {
                    const cells = row.split('|').filter(c => c.trim()).map(c => `<td>${c.trim()}</td>`).join('');
                    return `<tr>${cells}</tr>`;
                }).join('');
                return `<table><thead><tr>${headers}</tr></thead><tbody>${rows}</tbody></table>`;
            });

            // 标题（支持 # 到 ######）
            html = html.replace(/^###### (.*$)/gm, '<h6>$1</h6>');
            html = html.replace(/^##### (.*$)/gm, '<h5>$1</h5>');
            html = html.replace(/^#### (.*$)/gm, '<h4>$1</h4>');
            html = html.replace(/^### (.*$)/gm, '<h3>$1</h3>');
            html = html.replace(/^## (.*$)/gm, '<h2>$1</h2>');
            html = html.replace(/^# (.*$)/gm, '<h1>$1</h1>');

            // 任务列表
            html = html.replace(/^\- \[x\] (.*$)/gm, '<li class="task-item"><input type="checkbox" checked disabled> $1</li>');
            html = html.replace(/^\- \[ \] (.*$)/gm, '<li class="task-item"><input type="checkbox" disabled> $1</li>');

            // 粗体和斜体（支持嵌套）
            html = html.replace(/\*\*\*(.*?)\*\*\*/g, '<strong><em>$1</em></strong>');
            html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
            html = html.replace(/\*(.*?)\*/g, '<em>$1</em>');
            html = html.replace(/__(.*?)__/g, '<strong>$1</strong>');
            html = html.replace(/_(.*?)_/g, '<em>$1</em>');

            // 删除线
            html = html.replace(/~~(.*?)~~/g, '<del>$1</del>');

            // 高亮文本
            html = html.replace(/==(.*?)==/g, '<mark>$1</mark>');

            // 上标和下标
            html = html.replace(/\^(.*?)\^/g, '<sup>$1</sup>');
            html = html.replace(/~(.*?)~/g, '<sub>$1</sub>');

            // 图片（必须在链接之前处理）
            html = html.replace(/!\[(.*?)\]\((.*?) "(.*?)"\)/g, '<img src="$2" alt="$1" title="$3">');
            html = html.replace(/!\[(.*?)\]\((.*?)\)/g, '<img src="$2" alt="$1">');

            // 链接（支持标题）
            html = html.replace(/\[(.*?)\]\((.*?) "(.*?)"\)/g, '<a href="$2" title="$3">$1</a>');
            html = html.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2">$1</a>');

            // 自动链接（邮箱和URL）
            html = html.replace(/<(https?:\/\/[^>]+)>/g, '<a href="$1" target="_blank">$1</a>');
            html = html.replace(/<([^>]+@[^>]+)>/g, '<a href="mailto:$1">$1</a>');

            // 引用（支持多行）
            html = html.replace(/^> (.*$)/gm, '<blockquote>$1</blockquote>');
            html = html.replace(/<\/blockquote>\n<blockquote>/g, '\n');

            // 无序列表（支持嵌套）
            html = html.replace(/^[\-\*] (.*$)/gm, '<li>$1</li>');

            // 有序列表（支持嵌套）
            html = html.replace(/^\d+\. (.*$)/gm, '<li>$1</li>');

            // 水平线
            html = html.replace(/^---$/gm, '<hr>');
            html = html.replace(/^\*\*\*$/gm, '<hr>');
            html = html.replace(/^___$/gm, '<hr>');

            // 段落（双换行）
            html = html.replace(/\n\n/g, '</p><p>');

            // 单换行
            html = html.replace(/\n/g, '<br>');

            // 包装列表项
            html = html.replace(/(<li>.*?<\/li>)/gs, '<ul>$1</ul>');
            html = html.replace(/<\/ul>\s*<ul>/g, '');

            // 还原代码块和行内代码
            html = html.replace(/%%CODEBLOCK_(\d+)%%/g, (match, index) => codeBlocks[parseInt(index)]);
            html = html.replace(/%%INLINECODE_(\d+)%%/g, (match, index) => inlineCodes[parseInt(index)]);

            return '<p>' + html + '</p>';
        }

        // HTML 转义函数
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        // ========== 更新预览 ==========
        function updatePreview() {
            const content = document.getElementById('articleContent').value;
            const title = document.getElementById('articleTitle').value;
            const previewEl = document.getElementById('previewContent');

            if (!title && !content) {
                previewEl.innerHTML = `
                    <div class="preview-placeholder">
                        <div class="preview-placeholder-icon"><svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg></div>
                        <div>在左侧输入内容，这里将实时预览</div>
                    </div>
                `;
                return;
            }

            let html = '';
            if (title) {
                html += `<h1>${escapeHtml(title)}</h1>`;
            }
            html += markdownToHtml(content);
            previewEl.innerHTML = html;
        }

        // ========== 渲染文章卡片 ==========
        function renderArticleTable() {
            const container = document.getElementById('articleCardsContainer');
            const searchTerm = document.getElementById('articleSearch').value.toLowerCase();
            const categoryFilter = document.getElementById('categoryFilter').value;
            const statusFilter = document.getElementById('statusFilter').value;

            let filtered = state.articles.filter(article => {
                const matchSearch = article.title.toLowerCase().includes(searchTerm) ||
                                   article.summary.toLowerCase().includes(searchTerm);
                const matchCategory = !categoryFilter || article.category === categoryFilter;
                const matchStatus = !statusFilter || article.status === statusFilter;
                return matchSearch && matchCategory && matchStatus;
            });

            if (filtered.length === 0) {
                container.innerHTML = `
                    <div class="empty-state">
                        <div class="empty-icon"><svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline></svg></div>
                        <div class="empty-title">没有找到文章</div>
                        <div class="empty-desc">试试调整筛选条件或创建新文章</div>
                    </div>
                `;
                return;
            }

            container.innerHTML = filtered.map(article => `
                <article class="article-card" data-id="${article.id}">
                    <div class="article-card-header">
                        <span class="article-category">
                            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
                            ${escapeHtml(article.category)}
                        </span>
                        <span class="article-status ${article.status === 'published' ? 'status-published' : 'status-draft'}">
                            <span class="status-dot"></span>
                            ${article.status === 'published' ? '已发布' : '草稿'}
                        </span>
                    </div>
                    <div class="article-card-body">
                        <h3 class="article-title">${escapeHtml(article.title)}</h3>
                        <p class="article-summary">${escapeHtml(article.summary)}</p>
                        <div class="article-tags">
                            ${(article.tags || []).map(tag => `<span class="article-tag">${escapeHtml(tag)}</span>`).join('')}
                        </div>
                    </div>
                    <div class="article-card-footer">
                        <div class="article-meta">
                            <span class="meta-item">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"></rect><line x1="16" x2="16" y1="2" y2="6"></line><line x1="8" x2="8" y1="2" y2="6"></line><line x1="3" x2="21" y1="10" y2="10"></line></svg>
                                ${escapeHtml(article.createdAt)}
                            </span>
                            <span class="meta-item">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                                ${article.views.toLocaleString()}
                            </span>
                        </div>
                        <div class="article-actions">
                            <button class="action-btn btn-edit edit-btn" data-id="${article.id}">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                                编辑
                            </button>
                            <button class="action-btn btn-hide hide-btn" data-id="${article.id}" data-status="${article.status}">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                                ${article.status === 'published' ? '隐藏' : '显示'}
                            </button>
                            <button class="action-btn btn-delete delete-btn" data-id="${article.id}">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                                删除
                            </button>
                        </div>
                    </div>
                </article>
            `).join('');

            document.getElementById('totalArticles').textContent = filtered.length;
        }

        // ========== 更新统计 ==========
        function updateStats() {
            const total = state.articles.length;
            const published = state.articles.filter(a => a.status === 'published').length;
            const draft = total - published;
            const views = state.articles.reduce((sum, a) => sum + a.views, 0);

            document.getElementById('statTotal').textContent = total;
            document.getElementById('statPublished').textContent = published;
            document.getElementById('statDraft').textContent = draft;
            document.getElementById('statViews').textContent = views.toLocaleString();
        }

        // ========== 更新友链数量 ==========
        function updateLinkCount() {
            const linkCount = document.getElementById('linkCount');
            if (linkCount) {
                linkCount.textContent = state.links.length;
            }
        }

        // ========== 渲染友情链接表格 ==========
        function renderLinkTable() {
            const tbody = document.getElementById('linkTableBody');
            if (!tbody) return; // 卡片布局下无需渲染表格
            if (!state.links.length) {
                tbody.innerHTML = `
                    <tr>
                        <td colspan="4">
                            <div class="empty-state">
                                <div class="empty-icon"><svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path></svg></div>
                                <div class="empty-title">暂无友情链接</div>
                                <div class="empty-desc">点击上方按钮添加第一个链接</div>
                            </div>
                        </td>
                    </tr>
                `;
                return;
            }
            tbody.innerHTML = state.links.map(link => `
                <tr data-id="${link.id}">
                    <td style="font-weight: 500; color: var(--card-text-color-main);">${escapeHtml(link.name)}</td>
                    <td><a href="${escapeHtml(link.url)}" target="_blank" style="color: var(--accent-color);">${escapeHtml(link.url)}</a></td>
                    <td style="color: var(--card-text-color-secondary);">${escapeHtml(link.desc)}</td>
                    <td>
                        <div style="display: flex; gap: 8px;">
                            <button class="btn btn-secondary btn-sm edit-link-btn" data-id="${link.id}">编辑</button>
                            <button class="btn btn-danger btn-sm delete-link-btn" data-id="${link.id}">删除</button>
                        </div>
                    </td>
                </tr>
            `).join('');
        }

        // ========== 渲染每日一问表格 ==========
        function renderDqTable() {
            const tbody = document.getElementById('dqTableBody');
            if (!state.dailyQuestions.length) {
                tbody.innerHTML = `
                    <tr>
                        <td colspan="6">
                            <div class="empty-state">
                                <div class="empty-icon"><svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line></svg></div>
                                <div class="empty-title">暂无问题</div>
                                <div class="empty-desc">点击上方按钮创建第一个问题</div>
                            </div>
                        </td>
                    </tr>
                `;
                return;
            }
            tbody.innerHTML = state.dailyQuestions.map(q => `
                <tr data-id="${q.id}">
                    <td style="font-weight: 500; color: var(--card-text-color-main); max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" title="${escapeHtml(q.question)}">${escapeHtml(q.question)}</td>
                    <td style="color: var(--card-text-color-secondary);">${q.display_date ? new Date(q.display_date).toLocaleDateString() : '未设置'}</td>
                    <td>
                        <span class="status-badge ${q.status === 1 ? 'status-published' : 'status-draft'}">
                            <span class="status-dot"></span>
                            ${q.status === 1 ? '启用' : '禁用'}
                        </span>
                    </td>
                    <td style="color: var(--card-text-color-secondary);">
                        <span title="浏览">👁️ ${q.view_count}</span>
                        <span style="margin-left: 8px;" title="点赞">❤️ ${q.like_count}</span>
                    </td>
                    <td style="color: var(--card-text-color-tertiary);">${new Date(q.created_at).toLocaleDateString()}</td>
                    <td>
                        <div style="display: flex; gap: 8px;">
                            <button class="btn btn-secondary btn-sm edit-dq-btn" data-id="${q.id}">编辑</button>
                            <button class="btn btn-danger btn-sm delete-dq-btn" data-id="${q.id}">删除</button>
                        </div>
                    </td>
                </tr>
            `).join('');
        }

        // ========== 渲染每日一问评论表格 ==========
        function renderDqCommentTable() {
            const tbody = document.getElementById('dqCommentTableBody');
            if (!state.dailyQuestionComments.length) {
                tbody.innerHTML = `
                    <tr>
                        <td colspan="6">
                            <div class="empty-state">
                                <div class="empty-icon"><svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg></div>
                                <div class="empty-title">暂无评论</div>
                                <div class="empty-desc">等待用户发表评论</div>
                            </div>
                        </td>
                    </tr>
                `;
                return;
            }
            tbody.innerHTML = state.dailyQuestionComments.map(c => `
                <tr data-id="${c.id}">
                    <td style="font-weight: 500; color: var(--card-text-color-main); max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" title="${escapeHtml(c.question || '')}">${escapeHtml(c.question || '未知问题')}</td>
                    <td style="color: var(--card-text-color-secondary); max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" title="${escapeHtml(c.content)}">${escapeHtml(c.content)}</td>
                    <td style="color: var(--card-text-color-secondary);">${escapeHtml(c.nickname)}</td>
                    <td>
                        <span class="status-badge ${c.status === 1 ? 'status-published' : c.status === 2 ? 'status-draft' : 'status-pending'}">
                            <span class="status-dot"></span>
                            ${c.status === 0 ? '待审核' : c.status === 1 ? '已通过' : '已拒绝'}
                        </span>
                    </td>
                    <td style="color: var(--card-text-color-tertiary);">${new Date(c.created_at).toLocaleDateString()}</td>
                    <td>
                        <div style="display: flex; gap: 8px;">
                            ${c.status === 0 ? `
                                <button class="btn btn-secondary btn-sm approve-dq-comment-btn" data-id="${c.id}">通过</button>
                                <button class="btn btn-danger btn-sm reject-dq-comment-btn" data-id="${c.id}">拒绝</button>
                            ` : ''}
                            <button class="btn btn-danger btn-sm delete-dq-comment-btn" data-id="${c.id}">删除</button>
                        </div>
                    </td>
                </tr>
            `).join('');
        }

        // ========== 渲染每日一问统计 ==========
        function renderDqStats() {
            document.getElementById('dqStatTotal').textContent = state.dailyQuestionStats.total_questions;
            document.getElementById('dqStatEnabled').textContent = state.dailyQuestionStats.enabled_questions;
            document.getElementById('dqStatViews').textContent = state.dailyQuestionStats.total_views.toLocaleString();
            document.getElementById('dqStatLikes').textContent = state.dailyQuestionStats.total_likes.toLocaleString();
        }

        // ========== 渲染页面管理表格 ==========
        function renderPageTable() {
            const tbody = document.getElementById('pageTableBody');
            if (!state.pages.length) {
                tbody.innerHTML = `
                    <tr>
                        <td colspan="5">
                            <div class="empty-state">
                                <div class="empty-icon"><svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path></svg></div>
                                <div class="empty-title">暂无页面</div>
                                <div class="empty-desc">点击上方按钮创建第一个页面</div>
                            </div>
                        </td>
                    </tr>
                `;
                return;
            }
            tbody.innerHTML = state.pages.map(page => `
                <tr data-id="${page.id}">
                    <td style="font-weight: 500; color: var(--card-text-color-main);">${escapeHtml(page.title)}</td>
                    <td><code style="background: var(--body-background); padding: 2px 8px; border-radius: 4px; font-size: 13px;">/${escapeHtml(page.slug)}</code></td>
                    <td>
                        <span class="status-badge ${page.status === 'published' ? 'status-published' : 'status-draft'}">
                            <span class="status-dot"></span>
                            ${page.status === 'published' ? '已发布' : '草稿'}
                        </span>
                    </td>
                    <td style="color: var(--card-text-color-tertiary);">${page.updatedAt}</td>
                    <td>
                        <div style="display: flex; gap: 8px;">
                            <button class="btn btn-secondary btn-sm edit-page-btn" data-id="${page.id}">编辑</button>
                            <button class="btn btn-danger btn-sm delete-page-btn" data-id="${page.id}">删除</button>
                        </div>
                    </td>
                </tr>
            `).join('');
        }

        // ========== 渲染媒体库 ==========
        function renderMediaGrid() {
            const grid = document.getElementById('mediaGrid');
            if (!state.media.length) {
                grid.innerHTML = `
                    <div style="grid-column: 1 / -1; text-align: center; padding: 40px 20px; color: var(--card-text-color-tertiary);">
                        <div style="margin-bottom: 12px;"><svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg></div>
                        <div>暂无文件，点击上方按钮上传</div>
                    </div>
                `;
                return;
            }
            grid.innerHTML = state.media.map(item => `
                <div class="card" style="overflow: hidden; cursor: pointer;" data-id="${item.id}">
                    <div style="aspect-ratio: 1; background: var(--body-background); display: flex; align-items: center; justify-content: center; overflow: hidden;">
                        ${item.type.startsWith('image/')
                            ? `<img src="${item.data}" alt="${escapeHtml(item.name)}" style="width: 100%; height: 100%; object-fit: cover;">`
                            : `<div style="font-size: 40px;"><svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline></svg></div>`
                        }
                    </div>
                    <div style="padding: 10px 12px;">
                        <div style="font-size: 13px; font-weight: 500; color: var(--card-text-color-main); overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" title="${escapeHtml(item.name)}">${escapeHtml(item.name)}</div>
                        <div style="font-size: 11px; color: var(--card-text-color-tertiary); margin-top: 4px;">${item.size}</div>
                    </div>
                </div>
            `).join('');
        }

        // ========== 加载系统设置 ==========
        function loadSettings() {
            document.getElementById('siteName').value = state.settings.siteName || '';
            document.getElementById('siteDesc').value = state.settings.siteDesc || '';
            document.getElementById('siteUrl').value = state.settings.siteUrl || '';
            document.getElementById('pageSize').value = state.settings.pageSize || 10;
            document.getElementById('seoTitle').value = state.settings.seoTitle || '';
            document.getElementById('seoDesc').value = state.settings.seoDesc || '';
            document.getElementById('seoKeywords').value = state.settings.seoKeywords || '';
        }

        // ========== 保存系统设置 ==========
        function saveSettings() {
            state.settings = {
                siteName: document.getElementById('siteName').value.trim(),
                siteDesc: document.getElementById('siteDesc').value.trim(),
                siteUrl: document.getElementById('siteUrl').value.trim(),
                pageSize: parseInt(document.getElementById('pageSize').value) || 10,
                seoTitle: document.getElementById('seoTitle').value.trim(),
                seoDesc: document.getElementById('seoDesc').value.trim(),
                seoKeywords: document.getElementById('seoKeywords').value.trim()
            };
            localStorage.setItem(SETTINGS_KEY, JSON.stringify(state.settings));
            showToast('设置已保存', 'success');
        }

        // ========== 格式化文件大小 ==========
        function formatFileSize(bytes) {
            if (bytes < 1024) return bytes + ' B';
            if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
            return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
        }

        // ========== 切换页面 ==========
        function navigateTo(pageName) {
            switchPage(pageName);
        }

        function switchPage(pageName) {
            document.querySelectorAll('.menu-item').forEach(item => {
                if (item.dataset.page) {
                    item.classList.toggle('active', item.dataset.page === pageName);
                }
            });

            document.querySelectorAll('.page').forEach(page => {
                page.classList.toggle('active', page.id === `page-${pageName}`);
            });

            const titles = {
                home: '仪表盘',
                articles: '文章',
                links: '友链',
                entertainment: '娱乐',
                about: '关于我',
                edit: state.editingArticle ? '编辑文章' : '新建文章',
                categories: '分类管理',
                tags: '标签管理',
                comments: '评论管理',
                media: '媒体库',
                pages: '页面管理',
                stats: '数据统计',
                settings: '系统设置'
            };
            document.getElementById('pageTitle').textContent = titles[pageName] || pageName;

            state.currentPage = pageName;
        }

        // ========== 打开编辑页面 ==========
        function openEditPage(articleId = null) {
            if (articleId) {
                const article = state.articles.find(a => a.id === articleId);
                if (!article) return;

                state.editingArticle = article;
                document.getElementById('articleTitle').value = article.title;
                document.getElementById('articleContent').value = article.content;
            } else {
                state.editingArticle = null;
                document.getElementById('articleTitle').value = '';
                document.getElementById('articleContent').value = '';
            }

            switchPage('edit');
            updatePreview();
        }

        // ========== 保存文章 ==========
        function saveArticle(status, category, tags, summary) {
            const title = document.getElementById('articleTitle').value.trim();
            const content = document.getElementById('articleContent').value;

            if (!title) {
                showToast('请输入文章标题', 'error');
                return false;
            }
            if (!category) {
                showToast('请选择文章分类', 'error');
                return false;
            }

            const now = new Date().toISOString().split('T')[0];

            if (state.editingArticle) {
                const index = state.articles.findIndex(a => a.id === state.editingArticle.id);
                if (index !== -1) {
                    state.articles[index] = {
                        ...state.articles[index],
                        title,
                        content,
                        category,
                        tags: tags,
                        status: status,
                        summary: summary || content.substring(0, 200).replace(/[#*`>\-\[\]]/g, ''),
                        updatedAt: now
                    };
                }
                showToast('文章已更新', 'success');
            } else {
                const newArticle = {
                    id: Date.now(),
                    title,
                    content,
                    category,
                    tags: tags,
                    status: status,
                    summary: summary || content.substring(0, 200).replace(/[#*`>\-\[\]]/g, ''),
                    views: 0,
                    cover: null,
                    createdAt: now,
                    updatedAt: now
                };
                state.articles.unshift(newArticle);
                showToast('文章已创建', 'success');
            }

            saveToStorage();
            renderArticleTable();
            updateStats();
            switchPage('articles');
            return true;
        }

        // ========== 删除文章 ==========
        function deleteArticle(id) {
            state.deleteArticleId = id;
            document.getElementById('deleteModal').classList.add('active');
        }

        function confirmDeleteArticle() {
            if (state.deleteArticleId) {
                state.articles = state.articles.filter(a => a.id !== state.deleteArticleId);
                saveToStorage();
                renderArticleTable();
                updateStats();
                showToast('文章已删除', 'success');
            }
            document.getElementById('deleteModal').classList.remove('active');
            state.deleteArticleId = null;
        }

        // ========== 切换文章可见性 ==========
        function toggleArticleVisibility(id) {
            const article = state.articles.find(a => a.id === id);
            if (article) {
                article.status = article.status === 'published' ? 'draft' : 'published';
                article.updatedAt = new Date().toISOString().split('T')[0];
                saveToStorage();
                renderArticleTable();
                updateStats();
                showToast(article.status === 'published' ? '文章已显示' : '文章已隐藏', 'success');
            }
        }

        // ========== 标签管理（弹窗） ==========
        function renderModalTags() {
            const container = document.getElementById('modalTagsContainer');
            const input = document.getElementById('modalTagInput');

            container.innerHTML = '';
            state.modalTags.forEach(tag => {
                const tagEl = document.createElement('span');
                tagEl.className = 'tag-item';
                tagEl.innerHTML = `
                    ${escapeHtml(tag)}
                    <button class="tag-remove" data-tag="${escapeHtml(tag)}">×</button>
                `;
                container.appendChild(tagEl);
            });
            container.appendChild(input);
        }

        function addModalTag(tagName) {
            const tag = tagName.trim();
            if (tag && !state.modalTags.includes(tag)) {
                state.modalTags.push(tag);
                renderModalTags();
            }
        }

        function removeModalTag(tagName) {
            state.modalTags = state.modalTags.filter(t => t !== tagName);
            renderModalTags();
        }

        // ========== 打开发布弹窗 ==========
        function openPublishModal() {
            const title = document.getElementById('articleTitle').value.trim();
            if (!title) {
                showToast('请先输入文章标题', 'error');
                return;
            }

            // 如果是编辑已有文章，预填分类和标签
            if (state.editingArticle) {
                document.getElementById('modalCategory').value = state.editingArticle.category;
                document.getElementById('modalStatus').value = state.editingArticle.status;
                document.getElementById('modalSummary').value = state.editingArticle.summary;
                state.modalTags = [...state.editingArticle.tags];
            } else {
                document.getElementById('modalCategory').value = '';
                document.getElementById('modalStatus').value = 'published';
                document.getElementById('modalSummary').value = '';
                state.modalTags = [];
            }

            renderModalTags();
            document.getElementById('publishModal').classList.add('active');
        }

        function closePublishModal() {
            document.getElementById('publishModal').classList.remove('active');
        }

        // ========== 打开个人资料编辑弹窗 ==========
        function openProfileModal() {
            document.getElementById('profileUsername').value = state.profile.username;
            document.getElementById('profileDesc').value = state.profile.desc;

            // 渲染预览头像
            const preview = document.getElementById('profileAvatarPreview');
            if (state.profile.avatar) {
                preview.innerHTML = `<img src="${state.profile.avatar}" alt="头像"><div class="sidebar-avatar-edit">更换头像</div>`;
            } else {
                preview.innerHTML = `<span id="profileAvatarText">${state.profile.username.charAt(0).toUpperCase()}</span><div class="sidebar-avatar-edit">更换头像</div>`;
            }

            // 渲染社交链接列表
            renderSocialList();

            document.getElementById('profileModal').classList.add('active');
        }
        // 暴露到全局作用域，供 HTML onclick 使用
        window.openProfileModal = openProfileModal;

        function closeProfileModal() {
            document.getElementById('profileModal').classList.remove('active');
        }

        // ========== Toast 提示 ==========
        function showToast(message, type = 'success') {
            const container = document.getElementById('toastContainer');
            const toast = document.createElement('div');
            toast.className = `toast toast-${type}`;
            toast.innerHTML = `
                <span class="toast-icon">${type === 'success' ? '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>' : '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>'}</span>
                <span class="toast-message">${message}</span>
            `;
            container.appendChild(toast);

            setTimeout(() => {
                toast.remove();
            }, 3000);
        }

        // ========== HTML 转义 ==========
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        // ========== 每日一问辅助函数 ==========
        // 打开问题模态框
        function openDqModal(question) {
            const modal = document.createElement('div');
            modal.className = 'modal-overlay';
            modal.style.cssText = 'position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000;';

            const isEdit = question !== null;
            const title = isEdit ? '编辑问题' : '新建问题';

            modal.innerHTML = `
                <div class="card" style="width: 90%; max-width: 600px; max-height: 90vh; overflow-y: auto;">
                    <div class="card-header">
                        <div class="card-title">${title}</div>
                        <button class="btn btn-secondary btn-sm close-dq-modal-btn">✕</button>
                    </div>
                    <div class="card-body">
                        <div class="form-group">
                            <label class="form-label">问题 <span class="required">*</span></label>
                            <textarea class="form-textarea" id="dqQuestionInput" placeholder="输入问题内容..." style="min-height: 80px;">${isEdit ? escapeHtml(question.question) : ''}</textarea>
                        </div>
                        <div class="form-group">
                            <label class="form-label">答案 <span class="required">*</span></label>
                            <textarea class="form-textarea" id="dqAnswerInput" placeholder="输入答案内容..." style="min-height: 120px;">${isEdit ? escapeHtml(question.answer) : ''}</textarea>
                        </div>
                        <div class="form-group">
                            <label class="form-label">显示日期</label>
                            <input type="date" class="form-input" id="dqDateInput" value="${isEdit && question.display_date ? new Date(question.display_date).toISOString().split('T')[0] : ''}">
                            <div style="font-size: 12px; color: var(--card-text-color-tertiary); margin-top: 4px;">留空则按日期自动轮询显示</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">状态</label>
                            <select class="form-select" id="dqStatusInput">
                                <option value="1" ${isEdit && question.status === 1 ? 'selected' : ''}>启用</option>
                                <option value="0" ${isEdit && question.status === 0 ? 'selected' : ''}>禁用</option>
                            </select>
                        </div>
                    </div>
                    <div class="card-footer" style="display: flex; justify-content: flex-end; gap: 12px; padding: 16px; border-top: 1px solid var(--card-separator-color);">
                        <button class="btn btn-secondary close-dq-modal-btn">取消</button>
                        <button class="btn btn-primary save-dq-btn" data-id="${isEdit ? question.id : ''}">保存</button>
                    </div>
                </div>
            `;

            document.body.appendChild(modal);

            // 关闭模态框
            modal.querySelectorAll('.close-dq-modal-btn').forEach(btn => {
                btn.addEventListener('click', () => modal.remove());
            });

            // 点击遮罩关闭
            modal.addEventListener('click', (e) => {
                if (e.target === modal) modal.remove();
            });

            // 保存问题
            modal.querySelector('.save-dq-btn').addEventListener('click', () => {
                const questionText = document.getElementById('dqQuestionInput').value.trim();
                const answerText = document.getElementById('dqAnswerInput').value.trim();
                const displayDate = document.getElementById('dqDateInput').value;
                const status = parseInt(document.getElementById('dqStatusInput').value);

                if (!questionText) {
                    showToast('请输入问题内容', 'error');
                    return;
                }

                if (!answerText) {
                    showToast('请输入答案内容', 'error');
                    return;
                }

                const data = {
                    question: questionText,
                    answer: answerText,
                    display_date: displayDate || null,
                    status: status
                };

                if (isEdit) {
                    updateDq(question.id, data);
                } else {
                    createDq(data);
                }

                modal.remove();
            });
        }

        // 创建问题
        function createDq(data) {
            // 这里应该调用后端 API
            // 由于目前是前端演示，我们使用 localStorage 模拟
            const newQuestion = {
                id: Date.now(),
                ...data,
                view_count: 0,
                like_count: 0,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString()
            };

            state.dailyQuestions.unshift(newQuestion);
            localStorage.setItem('dailyQuestions', JSON.stringify(state.dailyQuestions));
            renderDqTable();
            updateDqStats();
            showToast('问题创建成功', 'success');
        }

        // 更新问题
        function updateDq(id, data) {
            const index = state.dailyQuestions.findIndex(q => q.id === id);
            if (index !== -1) {
                state.dailyQuestions[index] = {
                    ...state.dailyQuestions[index],
                    ...data,
                    updated_at: new Date().toISOString()
                };
                localStorage.setItem('dailyQuestions', JSON.stringify(state.dailyQuestions));
                renderDqTable();
                updateDqStats();
                showToast('问题更新成功', 'success');
            }
        }

        // 删除问题
        function deleteDq(id) {
            state.dailyQuestions = state.dailyQuestions.filter(q => q.id !== id);
            localStorage.setItem('dailyQuestions', JSON.stringify(state.dailyQuestions));
            renderDqTable();
            updateDqStats();
            showToast('问题已删除', 'success');
        }

        // 审核通过评论
        function approveDqComment(id) {
            const comment = state.dailyQuestionComments.find(c => c.id === id);
            if (comment) {
                comment.status = 1;
                localStorage.setItem('dailyQuestionComments', JSON.stringify(state.dailyQuestionComments));
                renderDqCommentTable();
                showToast('评论已通过', 'success');
            }
        }

        // 拒绝评论
        function rejectDqComment(id) {
            const comment = state.dailyQuestionComments.find(c => c.id === id);
            if (comment) {
                comment.status = 2;
                localStorage.setItem('dailyQuestionComments', JSON.stringify(state.dailyQuestionComments));
                renderDqCommentTable();
                showToast('评论已拒绝', 'success');
            }
        }

        // 删除评论
        function deleteDqComment(id) {
            state.dailyQuestionComments = state.dailyQuestionComments.filter(c => c.id !== id);
            localStorage.setItem('dailyQuestionComments', JSON.stringify(state.dailyQuestionComments));
            renderDqCommentTable();
            showToast('评论已删除', 'success');
        }

        // 更新每日一问统计
        function updateDqStats() {
            state.dailyQuestionStats = {
                total_questions: state.dailyQuestions.length,
                enabled_questions: state.dailyQuestions.filter(q => q.status === 1).length,
                total_views: state.dailyQuestions.reduce((sum, q) => sum + (q.view_count || 0), 0),
                total_likes: state.dailyQuestions.reduce((sum, q) => sum + (q.like_count || 0), 0),
                total_comments: state.dailyQuestionComments.length
            };
            renderDqStats();
        }

        // 加载每日一问数据
        function loadDqData() {
            const savedQuestions = localStorage.getItem('dailyQuestions');
            if (savedQuestions) {
                state.dailyQuestions = JSON.parse(savedQuestions);
            }

            const savedComments = localStorage.getItem('dailyQuestionComments');
            if (savedComments) {
                state.dailyQuestionComments = JSON.parse(savedComments);
            }

            updateDqStats();
            renderDqTable();
            renderDqCommentTable();
        }

        // ========== 事件绑定 ==========
        function bindEvents() {
            // 菜单点击
            document.querySelectorAll('.menu-item[data-page]').forEach(item => {
                item.addEventListener('click', () => {
                    const page = item.dataset.page;
                    if (page) {
                        switchPage(page);
                        document.getElementById('sidebar').classList.remove('open');
                    }
                });
            });

            // 移动端菜单
            document.getElementById('mobileMenuBtn').addEventListener('click', () => {
                document.getElementById('sidebar').classList.toggle('open');
            });

            // 主题切换
            document.getElementById('themeToggle').addEventListener('click', toggleTheme);

            // 语言切换
            document.getElementById('languageToggle').addEventListener('click', toggleLanguage);

            // 编辑个人资料
            document.getElementById('editProfileBtn').addEventListener('click', openProfileModal);
            document.getElementById('closeProfileModal').addEventListener('click', closeProfileModal);
            document.getElementById('cancelProfile').addEventListener('click', closeProfileModal);

            // 保存个人资料
            document.getElementById('saveProfile').addEventListener('click', () => {
                const username = document.getElementById('profileUsername').value.trim();
                const desc = document.getElementById('profileDesc').value.trim();

                if (!username) {
                    showToast('请输入用户名', 'error');
                    return;
                }

                // 收集社交链接数据
                const socialItems = document.querySelectorAll('.social-item');
                const social = [];
                socialItems.forEach(item => {
                    const iconInput = item.querySelector('.social-icon-input');
                    const urlInput = item.querySelector('.social-url-input');
                    if (iconInput && urlInput) {
                        social.push({
                            icon: iconInput.value.trim() || '●',
                            name: '',
                            url: urlInput.value.trim() || '#'
                        });
                    }
                });

                state.profile.username = username;
                state.profile.desc = desc;
                state.profile.social = social;
                saveProfile();
                renderProfile();
                closeProfileModal();
                showToast('资料已更新', 'success');
            });

            // 头像上传（侧边栏）
            document.getElementById('sidebarAvatar').addEventListener('click', () => {
                document.getElementById('avatarFileInput').click();
            });

            document.getElementById('avatarFileInput').addEventListener('change', (e) => {
                const file = e.target.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (event) => {
                        state.profile.avatar = event.target.result;
                        saveProfile();
                        renderProfile();
                        showToast('头像已更新', 'success');
                    };
                    reader.readAsDataURL(file);
                }
            });

            // 头像上传（弹窗预览）
            document.getElementById('profileAvatarPreview').addEventListener('click', () => {
                document.getElementById('profileAvatarInput').click();
            });

            document.getElementById('profileAvatarInput').addEventListener('change', (e) => {
                const file = e.target.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (event) => {
                        state.profile.avatar = event.target.result;
                        // 更新弹窗预览
                        const preview = document.getElementById('profileAvatarPreview');
                        preview.innerHTML = `<img src="${event.target.result}" alt="头像"><div class="sidebar-avatar-edit">更换头像</div>`;
                    };
                    reader.readAsDataURL(file);
                }
            });

            // 添加社交链接
            document.getElementById('addSocialBtn').addEventListener('click', () => {
                if (!state.profile.social) state.profile.social = [];
                state.profile.social.push({ icon: '●', name: '', url: '#' });
                renderSocialList();
            });

            // 删除社交链接
            document.getElementById('socialList').addEventListener('click', (e) => {
                const btn = e.target.closest('.remove-social-btn');
                if (btn) {
                    const index = parseInt(btn.dataset.index);
                    state.profile.social.splice(index, 1);
                    renderSocialList();
                }
            });

            // 新建文章
            document.getElementById('newArticleBtn').addEventListener('click', () => {
                openEditPage(null);
            });

            // 返回列表
            document.getElementById('backToListBtn').addEventListener('click', () => {
                switchPage('articles');
            });

            // 顶部返回列表按钮
            document.getElementById('backToListBtnTop').addEventListener('click', () => {
                switchPage('articles');
            });

            // 文章卡片操作
            document.getElementById('articleCardsContainer').addEventListener('click', (e) => {
                const editBtn = e.target.closest('.edit-btn');
                const deleteBtn = e.target.closest('.delete-btn');
                const hideBtn = e.target.closest('.hide-btn');

                if (editBtn) {
                    const id = parseInt(editBtn.dataset.id);
                    openEditPage(id);
                } else if (deleteBtn) {
                    const id = parseInt(deleteBtn.dataset.id);
                    deleteArticle(id);
                } else if (hideBtn) {
                    const id = parseInt(hideBtn.dataset.id);
                    toggleArticleVisibility(id);
                }
            });

            // 搜索和筛选
            document.getElementById('articleSearch').addEventListener('input', renderArticleTable);
            document.getElementById('categoryFilter').addEventListener('change', renderArticleTable);
            document.getElementById('statusFilter').addEventListener('change', renderArticleTable);

            // 编辑器实时预览
            document.getElementById('articleTitle').addEventListener('input', updatePreview);
            document.getElementById('articleContent').addEventListener('input', updatePreview);

            // 工具栏按钮功能
            document.querySelector('.editor-toolbar').addEventListener('click', (e) => {
                const btn = e.target.closest('.toolbar-btn');
                if (!btn) return;

                const action = btn.dataset.action;
                const textarea = document.getElementById('articleContent');
                const start = textarea.selectionStart;
                const end = textarea.selectionEnd;
                const selected = textarea.value.substring(start, end);

                let before = textarea.value.substring(0, start);
                let after = textarea.value.substring(end);
                let insertion = '';
                let cursorOffset = 0;

                switch (action) {
                    case 'bold':
                        insertion = `**${selected || '粗体文本'}**`;
                        cursorOffset = selected ? insertion.length : 2;
                        break;
                    case 'italic':
                        insertion = `*${selected || '斜体文本'}*`;
                        cursorOffset = selected ? insertion.length : 1;
                        break;
                    case 'strikethrough':
                        insertion = `~~${selected || '删除线文本'}~~`;
                        cursorOffset = selected ? insertion.length : 2;
                        break;
                    case 'highlight':
                        insertion = `==${selected || '高亮文本'}==`;
                        cursorOffset = selected ? insertion.length : 2;
                        break;
                    case 'h1':
                        insertion = `# ${selected || '标题1'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'h2':
                        insertion = `## ${selected || '标题2'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'h3':
                        insertion = `### ${selected || '标题3'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'ul':
                        insertion = `\n- ${selected || '列表项'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'ol':
                        insertion = `\n1. ${selected || '列表项'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'task':
                        insertion = `\n- [ ] ${selected || '任务项'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'quote':
                        insertion = `\n> ${selected || '引用文本'}`;
                        cursorOffset = insertion.length;
                        break;
                    case 'link':
                        insertion = `[${selected || '链接文本'}](url)`;
                        cursorOffset = selected ? insertion.length - 1 : 1;
                        break;
                    case 'image':
                        insertion = `![${selected || '图片描述'}](image-url)`;
                        cursorOffset = selected ? insertion.length - 1 : 2;
                        break;
                    case 'code':
                        if (selected.includes('\n')) {
                            insertion = `\`\`\`\n${selected}\n\`\`\``;
                        } else {
                            insertion = `\`${selected || '代码'}\``;
                        }
                        cursorOffset = selected ? insertion.length : 1;
                        break;
                    case 'table':
                        insertion = `\n| 列1 | 列2 | 列3 |\n| --- | --- | --- |\n| 内容 | 内容 | 内容 |`;
                        cursorOffset = insertion.length;
                        break;
                    case 'hr':
                        insertion = `\n\n---\n\n`;
                        cursorOffset = insertion.length;
                        break;
                }

                textarea.value = before + insertion + after;
                textarea.focus();
                textarea.selectionStart = start + cursorOffset;
                textarea.selectionEnd = start + cursorOffset;
                updatePreview();
            });

            // 保存草稿
            document.getElementById('saveDraftBtn').addEventListener('click', () => {
                saveArticle('draft', '搭建网站', [], '');
            });

            // 发布文章 - 打开弹窗
            document.getElementById('publishBtn').addEventListener('click', openPublishModal);

            // 发布弹窗 - 标签输入
            document.getElementById('modalTagInput').addEventListener('keydown', (e) => {
                if (e.key === 'Enter') {
                    e.preventDefault();
                    addModalTag(e.target.value);
                    e.target.value = '';
                }
            });

            // 发布弹窗 - 标签删除
            document.getElementById('modalTagsContainer').addEventListener('click', (e) => {
                if (e.target.classList.contains('tag-remove')) {
                    removeModalTag(e.target.dataset.tag);
                }
            });

            // 发布弹窗 - 封面图上传
            document.getElementById('modalCoverUpload').addEventListener('click', () => {
                document.getElementById('modalCoverFileInput').click();
            });

            // 发布弹窗 - 确认发布
            document.getElementById('confirmPublish').addEventListener('click', () => {
                const category = document.getElementById('modalCategory').value;
                const status = document.getElementById('modalStatus').value;
                const summary = document.getElementById('modalSummary').value.trim();

                if (!category) {
                    showToast('请选择文章分类', 'error');
                    return;
                }

                const success = saveArticle(status, category, state.modalTags, summary);
                if (success) {
                    closePublishModal();
                }
            });

            // 发布弹窗 - 取消
            document.getElementById('cancelPublish').addEventListener('click', closePublishModal);
            document.getElementById('closePublishModal').addEventListener('click', closePublishModal);

            // 模态框
            document.getElementById('closeDeleteModal').addEventListener('click', () => {
                document.getElementById('deleteModal').classList.remove('active');
            });
            document.getElementById('cancelDelete').addEventListener('click', () => {
                document.getElementById('deleteModal').classList.remove('active');
            });
            document.getElementById('confirmDelete').addEventListener('click', confirmDeleteArticle);

            // ========== 友情链接 ==========
            document.getElementById('newLinkBtn').addEventListener('click', () => {
                const name = prompt('网站名称：');
                if (!name) return;
                const url = prompt('链接地址：', 'https://');
                if (!url) return;
                const desc = prompt('网站描述：', '') || '';
                state.links.push({ id: Date.now(), name, url, desc });
                localStorage.setItem(LINKS_KEY, JSON.stringify(state.links));
                renderLinkTable();
                updateLinkCount();
                showToast('链接已添加', 'success');
            });

            const linkTableBody = document.getElementById('linkTableBody');
            if (linkTableBody) linkTableBody.addEventListener('click', (e) => {
                const editBtn = e.target.closest('.edit-link-btn');
                const deleteBtn = e.target.closest('.delete-link-btn');
                if (editBtn) {
                    const id = parseInt(editBtn.dataset.id);
                    const link = state.links.find(l => l.id === id);
                    if (link) {
                        const newName = prompt('网站名称：', link.name);
                        if (newName !== null) link.name = newName;
                        const newUrl = prompt('链接地址：', link.url);
                        if (newUrl !== null) link.url = newUrl;
                        const newDesc = prompt('网站描述：', link.desc);
                        if (newDesc !== null) link.desc = newDesc;
                        localStorage.setItem(LINKS_KEY, JSON.stringify(state.links));
                        renderLinkTable();
                        showToast('链接已更新', 'success');
                    }
                } else if (deleteBtn) {
                    if (confirm('确定删除此链接？')) {
                        const id = parseInt(deleteBtn.dataset.id);
                        state.links = state.links.filter(l => l.id !== id);
                        localStorage.setItem(LINKS_KEY, JSON.stringify(state.links));
                        renderLinkTable();
                        updateLinkCount();
                        showToast('链接已删除', 'success');
                    }
                }
            });

            // ========== 媒体库 ==========
            document.getElementById('uploadMediaBtn').addEventListener('click', () => {
                document.getElementById('mediaFileInput').click();
            });

            document.getElementById('mediaUploadArea').addEventListener('click', () => {
                document.getElementById('mediaFileInput').click();
            });

            document.getElementById('mediaFileInput').addEventListener('change', (e) => {
                const files = Array.from(e.target.files);
                files.forEach(file => {
                    if (file.size > 5 * 1024 * 1024) {
                        showToast(`${file.name} 超过 5MB 限制`, 'error');
                        return;
                    }
                    const reader = new FileReader();
                    reader.onload = (event) => {
                        state.media.push({
                            id: Date.now() + Math.random(),
                            name: file.name,
                            type: file.type,
                            size: formatFileSize(file.size),
                            data: event.target.result
                        });
                        localStorage.setItem(MEDIA_KEY, JSON.stringify(state.media));
                        renderMediaGrid();
                    };
                    reader.readAsDataURL(file);
                });
                showToast(`${files.length} 个文件已上传`, 'success');
                e.target.value = '';
            });

            // 拖拽上传
            const uploadArea = document.getElementById('mediaUploadArea');
            uploadArea.addEventListener('dragover', (e) => {
                e.preventDefault();
                uploadArea.style.borderColor = 'var(--accent-color)';
                uploadArea.style.background = 'rgba(var(--accent-color-rgb), 0.04)';
            });
            uploadArea.addEventListener('dragleave', () => {
                uploadArea.style.borderColor = '';
                uploadArea.style.background = '';
            });
            uploadArea.addEventListener('drop', (e) => {
                e.preventDefault();
                uploadArea.style.borderColor = '';
                uploadArea.style.background = '';
                const files = Array.from(e.dataTransfer.files);
                files.forEach(file => {
                    if (file.size > 5 * 1024 * 1024) {
                        showToast(`${file.name} 超过 5MB 限制`, 'error');
                        return;
                    }
                    const reader = new FileReader();
                    reader.onload = (event) => {
                        state.media.push({
                            id: Date.now() + Math.random(),
                            name: file.name,
                            type: file.type,
                            size: formatFileSize(file.size),
                            data: event.target.result
                        });
                        localStorage.setItem(MEDIA_KEY, JSON.stringify(state.media));
                        renderMediaGrid();
                    };
                    reader.readAsDataURL(file);
                });
                showToast(`${files.length} 个文件已上传`, 'success');
            });

            // ========== 页面管理 ==========
            document.getElementById('newPageBtn').addEventListener('click', () => {
                const title = prompt('页面标题：');
                if (!title) return;
                const slug = prompt('页面别名（URL 路径）：', title.toLowerCase().replace(/\s+/g, '-'));
                if (!slug) return;
                state.pages.push({ id: Date.now(), title, slug, status: 'draft', updatedAt: new Date().toISOString().split('T')[0] });
                localStorage.setItem(PAGES_KEY, JSON.stringify(state.pages));
                renderPageTable();
                showToast('页面已创建', 'success');
            });

            document.getElementById('pageTableBody').addEventListener('click', (e) => {
                const editBtn = e.target.closest('.edit-page-btn');
                const deleteBtn = e.target.closest('.delete-page-btn');
                if (editBtn) {
                    const id = parseInt(editBtn.dataset.id);
                    const page = state.pages.find(p => p.id === id);
                    if (page) {
                        const newTitle = prompt('页面标题：', page.title);
                        if (newTitle !== null) page.title = newTitle;
                        const newSlug = prompt('页面别名：', page.slug);
                        if (newSlug !== null) page.slug = newSlug;
                        const newStatus = confirm('设为已发布？（确定=已发布，取消=草稿）') ? 'published' : 'draft';
                        page.status = newStatus;
                        page.updatedAt = new Date().toISOString().split('T')[0];
                        localStorage.setItem(PAGES_KEY, JSON.stringify(state.pages));
                        renderPageTable();
                        showToast('页面已更新', 'success');
                    }
                } else if (deleteBtn) {
                    if (confirm('确定删除此页面？')) {
                        const id = parseInt(deleteBtn.dataset.id);
                        state.pages = state.pages.filter(p => p.id !== id);
                        localStorage.setItem(PAGES_KEY, JSON.stringify(state.pages));
                        renderPageTable();
                        showToast('页面已删除', 'success');
                    }
                }
            });

            // ========== 系统设置 ==========
            document.getElementById('saveSettingsBtn').addEventListener('click', saveSettings);

            // ========== 每日一问管理 ==========
            // 新建问题
            document.getElementById('newDqBtn').addEventListener('click', () => {
                openDqModal(null);
            });

            // 搜索和筛选
            document.getElementById('dqSearch').addEventListener('input', renderDqTable);
            document.getElementById('dqStatusFilter').addEventListener('change', renderDqTable);
            document.getElementById('dqCommentFilter').addEventListener('change', renderDqCommentTable);

            // 问题表格操作
            document.getElementById('dqTableBody').addEventListener('click', (e) => {
                const editBtn = e.target.closest('.edit-dq-btn');
                const deleteBtn = e.target.closest('.delete-dq-btn');

                if (editBtn) {
                    const id = parseInt(editBtn.dataset.id);
                    const question = state.dailyQuestions.find(q => q.id === id);
                    if (question) {
                        openDqModal(question);
                    }
                } else if (deleteBtn) {
                    const id = parseInt(deleteBtn.dataset.id);
                    if (confirm('确定删除此问题？')) {
                        deleteDq(id);
                    }
                }
            });

            // 评论表格操作
            document.getElementById('dqCommentTableBody').addEventListener('click', (e) => {
                const approveBtn = e.target.closest('.approve-dq-comment-btn');
                const rejectBtn = e.target.closest('.reject-dq-comment-btn');
                const deleteBtn = e.target.closest('.delete-dq-comment-btn');

                if (approveBtn) {
                    const id = parseInt(approveBtn.dataset.id);
                    approveDqComment(id);
                } else if (rejectBtn) {
                    const id = parseInt(rejectBtn.dataset.id);
                    rejectDqComment(id);
                } else if (deleteBtn) {
                    const id = parseInt(deleteBtn.dataset.id);
                    if (confirm('确定删除此评论？')) {
                        deleteDqComment(id);
                    }
                }
            });

            // ========== 卡片按钮事件委托 ==========

            // 友链卡片
            const linkCards = document.getElementById('linkCards');
            if (linkCards) {
                linkCards.addEventListener('click', (e) => {
                    const card = e.target.closest('.link-card');
                    if (e.target.closest('.btn-edit')) {
                        const name = card.querySelector('.link-card-name')?.textContent;
                        showToast(`编辑友链：${name}`, 'info');
                    } else if (e.target.closest('.btn-delete')) {
                        const name = card.querySelector('.link-card-name')?.textContent;
                        if (confirm(`确定删除友链「${name}」？`)) {
                            card.remove();
                            showToast('友链已删除', 'success');
                        }
                    }
                });
            }

            // 分类卡片
            document.querySelectorAll('.category-card').forEach(card => {
                card.addEventListener('click', (e) => {
                    if (e.target.closest('.btn-edit')) {
                        const name = card.querySelector('.category-card-title')?.textContent;
                        showToast(`编辑分类：${name}`, 'info');
                    } else if (e.target.closest('.btn-delete')) {
                        const name = card.querySelector('.category-card-title')?.textContent;
                        if (confirm(`确定删除分类「${name}」？`)) {
                            card.remove();
                            showToast('分类已删除', 'success');
                        }
                    }
                });
            });

            // 标签卡片
            document.querySelectorAll('.tag-card').forEach(card => {
                card.addEventListener('click', (e) => {
                    if (e.target.closest('.btn-edit')) {
                        const name = card.querySelector('.tag-card-name')?.textContent;
                        showToast(`编辑标签：${name}`, 'info');
                    } else if (e.target.closest('.btn-delete')) {
                        const name = card.querySelector('.tag-card-name')?.textContent;
                        if (confirm(`确定删除标签「${name}」？`)) {
                            card.remove();
                            showToast('标签已删除', 'success');
                        }
                    }
                });
            });

            // 评论卡片
            document.querySelectorAll('.comment-card').forEach(card => {
                card.addEventListener('click', (e) => {
                    if (e.target.closest('.btn-edit')) {
                        const name = card.querySelector('.comment-user-name')?.textContent;
                        showToast(`已通过 ${name} 的评论`, 'success');
                    } else if (e.target.closest('.btn-hide')) {
                        const name = card.querySelector('.comment-user-name')?.textContent;
                        showToast(`已拒绝 ${name} 的评论`, 'info');
                    } else if (e.target.closest('.btn-delete')) {
                        const name = card.querySelector('.comment-user-name')?.textContent;
                        if (confirm(`确定删除 ${name} 的评论？`)) {
                            card.remove();
                            showToast('评论已删除', 'success');
                        }
                    }
                });
            });

            // 娱乐卡片
            document.querySelectorAll('.media-card').forEach(card => {
                card.addEventListener('click', (e) => {
                    if (e.target.closest('.btn-edit')) {
                        const name = card.querySelector('.media-card-title')?.textContent;
                        showToast(`编辑：${name}`, 'info');
                    } else if (e.target.closest('.btn-delete')) {
                        const name = card.querySelector('.media-card-title')?.textContent;
                        if (confirm(`确定删除「${name}」？`)) {
                            card.remove();
                            showToast('已删除', 'success');
                        }
                    }
                });
            });

        }

        // 启动应用
        init();
