const articles = [
  { title: "My Blog Theme is Open Sourced, Welcome to Use It", description: "A ready-to-use blog template based on Hugo Theme Stack.", href: "#post-open-source-theme" },
  { title: "DesktopSnap: Your Open-Source Solution for Persistent Windows Desktop Layouts", description: "A portable WinUI3 tool to backup and restore your Windows icon layout.", href: "#post-desktopsnap" },
  { title: "Adding a Free AI Assistant with a Knowledge Base to My Blog using Cloudflare Workers", description: "Give your blog an AI assistant that answers questions based on your content.", href: "#post-cloudflare-ai" },
  { title: "Home Network Upgrade Plan: Upgrading to 2.5G LAN", description: "Upgrade your home network to 2.5G and install Gree Cloud Control.", href: "#post-net-upgrade" },
  { title: "Recording a QNAP NAS System Disk Migration Process", description: "How to migrate your QNAP NAS system disk.", href: "#post-nas-migra" },
  { title: "Purchase UPS Power Supply to Protect NAS and Computer", description: "Protect your NAS and computer with a UPS.", href: "#post-ups-nas" },
];

const postDB = {
  "post-open-source-theme": {
    title: "My Blog Theme is Open Sourced, Welcome to Use It",
    date: "2026-04-24",
    category: "Web",
    tags: ["Hugo", "Open Source", "Blog Template"],
    readTime: "3 minute read",
    content: '<p>Previously, some comments asked if I could open-source my blog modifications. However, because the theme was upgraded to V4 and I made a lot of changes, the process of upgrading and open-sourcing was quite troublesome.</p><p>So it was delayed for a long time, but it\'s finally online. Welcome to use it~ The open-source version is modified based on the Hugo-Theme-Stack V4 version.</p><p>This blog theme itself is very popular, and many people are using it. The official version is already great, and you can also use it directly.</p><p>Project URL: 👉 <a href="https://github.com/fcy222fcy/hugo-stack-starter" target="_blank"><strong>Hugo Stack Starter</strong></a></p><h2>Why build this template?</h2><p>I really like the minimalism and elegance of the <a href="https://github.com/CaiJimmy/hugo-theme-stack" target="_blank">Hugo Theme Stack</a>, which is highly popular on GitHub. Based on my personal preferences, I made a series of modifications:</p><ul><li>Added Mac-style three-color dots to code blocks and implemented automatic folding for long code snippets.</li><li>Refactored the mobile navigation bar, adding animation effects and a dark mode toggle.</li><li>Developed a complete /stats/ page to display my writing heat map and category habits.</li><li>Integrated the Waline comment system and beautified its styling.</li><li>Designed a two-column grid layout for the PC homepage.</li></ul><h2>What are its features?</h2><h3>🚀 Out of the Box</h3><p>I have set it as a <strong>GitHub Template</strong> and built in GitHub Actions workflows. You just need to click <code>Use this template</code> in your browser, change the domain name in <code>config.toml</code>, and your blog can be instantly deployed to GitHub Pages!</p><h3>🎨 Stylish Color Appearance</h3><p>The entire UI uses an eye-friendly color palette of <strong>Deep Sea Blue + Warm Beige</strong>. The dark mode has also been carefully tuned (using a black-grey pairing similar to GitHub Dark), and delicate micro-interaction animations have been added to areas like card hovering and menu expanding.</p><h3>📖 Built-in Tutorials</h3><p>To lower the learning curve for everyone, I have included 6 tutorial articles in the template. When you clone and run this blog, the homepage will no longer show empty placeholders.</p><h2>Conclusion</h2><p>If you also want to build your own blog or try out a new theme, why not give <a href="https://github.com/fcy222fcy/hugo-stack-starter" target="_blank">Hugo Stack Starter</a> a try! If you have any suggestions, feel free to open an Issue in the repository or just tell me in the comments section.</p>'
  },
  "post-desktopsnap": {
    title: "DesktopSnap: Your Open-Source Solution for Persistent Windows Desktop Layouts",
    date: "2026-03-12",
    category: "Software",
    tags: ["Windows", "Open Source", "WinUI3"],
    readTime: "3 minute read",
    content: '<p>Like many of you, I often encounter the frustration of Windows desktop icons getting disordered after switching resolutions, plugging/unplugging monitors, or after system updates. For multi-monitor users, this is especially painful — having to reposition all icons every time.</p><p>While there are some tools available online that can backup desktop layouts, most are bloated, have outdated interfaces, or require payment. So I decided to build my own.</p><h2>What is DesktopSnap?</h2><p>DesktopSnap is a lightweight desktop tool built with WinUI3. Its core functionality is simple: backup and restore desktop icon layouts. Clean interface, simple operation, one click done.</p><h2>Features</h2><ul><li><strong>One-click Backup</strong> — Read all icon positions on the desktop and save to a config file</li><li><strong>One-click Restore</strong> — Read config file and restore all icons to previous positions</li><li><strong>Multi-monitor Support</strong> — Correctly identifies icon positions across monitors</li><li><strong>Auto-restore on Boot</strong> — Optionally restore layout when Windows starts</li><li><strong>Multiple Backups</strong> — Save different layout profiles and switch between them</li></ul><h2>Technical Details</h2><p>Desktop icon position data is stored in the Windows Registry at: <code>HKCU\\Software\\Microsoft\\Windows\\Shell\\Bags\\1\\Desktop</code></p><p>By reading and writing FFLAGS, IconLayouts and other values, we can get and set icon positions. The actual implementation uses the IFolderView2 COM interface to manipulate Shell views.</p><h2>Download</h2><p>Available on Microsoft Store — search for "DesktopSnap" to download for free.</p>'
  },
  "post-cloudflare-ai": {
    title: "Adding a Free AI Assistant with a Knowledge Base to My Blog using Cloudflare Workers",
    date: "2026-04-21",
    category: "Web",
    tags: ["Cloudflare", "AI", "Workers"],
    readTime: "6 minute read",
    content: '<p>If you visit my blog using a desktop browser, you\'ll probably notice a tiny, blinking robot icon in the bottom-right corner. Click it and a chat panel slides out — you can ask any question about my blog content, and the AI will answer based on the articles I\'ve written.</p><p>This is my blog\'s AI assistant, built entirely using Cloudflare\'s free-tier infrastructure. No server needed, no API costs.</p><h2>Architecture Overview</h2><p>The solution has three parts:</p><ol><li><strong>Data Indexing</strong> — Split blog posts into chunks, generate embeddings, store in Cloudflare Vectorize</li><li><strong>Query Service</strong> — Convert questions to vectors, search for relevant article snippets</li><li><strong>AI Response</strong> — Send retrieved context + user question to an AI model</li></ol><h2>How It Works</h2><p>When a user asks a question, the system converts it to a vector embedding, searches Vectorize for the most relevant blog post chunks, then sends those chunks along with the question to Cloudflare Workers AI (currently using Llama 3) to generate a natural language answer.</p><h2>Cost</h2><p>Cloudflare Workers free tier: 100,000 requests/day. Vectorize free tier: 1,000,000 queries/day. For a personal blog, this is more than enough.</p>'
  },
  "post-net-upgrade": {
    title: "Home Network Upgrade Plan: Upgrading to 2.5G LAN and Installing Gree Cloud Control",
    date: "2025-07-29",
    category: "Knowledge Share",
    tags: ["Network"],
    readTime: "8 minute read",
    content: '<p>Recently, our home got a new broadband connection, so I decided to upgrade the internal network. The千兆 (gigabit) network was becoming a bottleneck when transferring large files between my NAS and computer.</p><h2>Upgrade Plan</h2><ul><li>Replace the switch with a 2.5G model</li><li>Add 2.5G network cards to NAS and computer</li><li>Use Cat6a cables</li></ul><h2>Gree Cloud Control</h2><p>While I was at it, I also installed a Gree Cloud Control module to enable remote control and smart features for the air conditioner.</p>'
  },
  "post-nas-migra": {
    title: "Recording a QNAP NAS System Disk Migration Process",
    date: "2025-07-21",
    category: "Knowledge Share",
    tags: ["NAS", "QNAP"],
    readTime: "6 minute read",
    content: '<p>My previous NAS system disk was a mechanical hard drive. After years of use, it was getting slower and slower, and I was worried about potential failures. Decided to migrate to an SSD.</p><h2>Why Migrate</h2><ul><li>Mechanical drive speeds are too slow for system operations</li><li>Risk of data loss from drive failure</li><li>SSD prices have dropped significantly</li></ul><h2>Migration Process</h2><p>The process involves backing up the current system, installing the new SSD, and restoring the system. QNAP provides built-in tools for this.</p>'
  },
  "post-ups-nas": {
    title: "Purchase UPS Power Supply to Protect NAS and Computer Safe Shutdown",
    date: "2025-07-18",
    category: "Knowledge Share",
    tags: ["NAS", "QNAP"],
    readTime: "5 minute read",
    content: '<p>With the summer heat approaching, in order to protect my NAS and Mac mini from unexpected power outages, I decided to purchase a UPS (Uninterruptible Power Supply).</p><h2>Why a UPS?</h2><p>Power outages can cause data corruption on NAS drives, especially during write operations. A UPS provides enough battery time to safely shut down all devices.</p><h2>What I Bought</h2><p>Selected a UPS with enough capacity to run my NAS and router for about 15 minutes — enough time for an automatic safe shutdown.</p>'
  },
  "post-komari": {
    title: "Sharing a Multi-Server Monitoring Tool - Komari",
    date: "2025-10-24",
    category: "Knowledge Share",
    tags: ["NAS", "Server"],
    readTime: "4 minute read",
    content: '<p>I\'ve always needed a tool to monitor multiple servers. I have several: one running my blog, one running Jellyfin, and one for development. The worst thing is when a server goes down and you don\'t know until someone tells you.</p><p>I\'ve tried Zabbix, Prometheus + Grafana, but they\'re overkill for personal use. Then I found Komari — a lightweight multi-server monitoring tool.</p><h2>Features</h2><ul><li>Real-time monitoring: CPU, memory, disk, network</li><li>Multi-server support: one dashboard for all servers</li><li>Alert notifications via email and webhooks</li><li>7-day historical data with trend charts</li><li>Lightweight deployment: single Docker container</li></ul><h2>Deployment</h2><p>One command to deploy: <code>docker run -d --name komari -p 8080:8080 komari/komari</code></p>'
  },
  "post-metatube": {
    title: "Free Metatube Backend Service Setup with Huggingface Space",
    date: "2025-09-12",
    category: "Knowledge Share",
    tags: ["Jellyfin", "Metatube"],
    readTime: "2 minute read",
    content: '<p>Huggingface has launched the Space feature, providing a free runtime environment. I used it to deploy a Metatube backend service for Jellyfin media scraping.</p><p>Metatube aggregates metadata from multiple sources (Douban, TMDB, Bilibili etc.) providing better Chinese metadata than TMDB alone.</p><h2>Deployment Steps</h2><ol><li>Register a Huggingface account</li><li>Create a new Space with Docker template</li><li>Write a Dockerfile</li><li>Configure environment variables</li><li>Deploy and get the access URL</li></ol>'
  },
  "post-game-record": {
    title: "Adding a Game Records Section to Hugo Blog with Automatic Game Data Fetching",
    date: "2025-09-03",
    category: "Knowledge Share",
    tags: ["Website", "Hugo"],
    readTime: "11 minute read",
    content: '<p>I\'ve always wanted to showcase the games I\'ve played on my blog — playtime, achievements, ratings. But manually entering this data is tedious and won\'t auto-update.</p><p>So I came up with a solution: use APIs to automatically fetch game data, sync it to Hugo data files, and render pages during build.</p><h2>Implementation</h2><ul><li>Steam API for game list and playtime</li><li>IGDB API for game covers and ratings</li><li>GitHub Actions for daily automated sync</li><li>Hugo templates for rendering</li></ul>'
  },
  "post-nextjs-adsense": {
    title: "Integrating Adsense into Next.js Website and Manually Placing Ads",
    date: "2025-09-02",
    category: "Knowledge Share",
    tags: ["Adsense", "Nextjs"],
    readTime: "6 minute read",
    content: '<p>In my previous article, I discussed the first issue I encountered when integrating Adsense — the preview not working on Cloudflare-deployed sites. In this article, I\'ll cover the actual integration and manual ad placement.</p><h2>Integration Steps</h2><ol><li>Register a Google AdSense account</li><li>Add your site and wait for approval</li><li>Add the AdSense script to your Next.js project</li><li>Manually place ad units</li></ol>'
  },
  "post-adsense-verify": {
    title: "Finally Verified AdSense Address and Received First Payment",
    date: "2024-12-27",
    category: "Knowledge Share",
    tags: ["Adsense", "Google"],
    readTime: "3 minute read",
    content: '<p>After nearly four months of waiting, I finally completed Google AdSense address verification and received my first payment!</p><h2>The Process</h2><ol><li>Receive the PIN code postcard</li><li>Enter the PIN code</li><li>Link a bank account</li><li>Wait for earnings to reach the minimum payout threshold</li><li>Receive the first bank transfer</li></ol><p>Though the amount isn\'t large, it\'s a validation of my content creation efforts.</p>'
  },
  "post-cloudflare-tunnel": {
    title: "Easy Intranet Access with Cloudflare Tunnel",
    date: "2024-09-27",
    category: "Knowledge Share",
    tags: ["Cloudflare"],
    readTime: "5 minute read",
    content: '<p>Cloudflare, the champion of free-tier users, offers an arsenal of powerful tools. One of them is Cloudflare Tunnel — which lets you expose services on your internal network to the internet without opening ports or having a public IP.</p><h2>How It Works</h2><p>You install a lightweight agent (cloudflared) on your server. It creates an outbound connection to Cloudflare\'s edge network. Then you configure a domain to route through the tunnel.</p><h2>Setup Steps</h2><ol><li>Install cloudflared</li><li>Login to Cloudflare</li><li>Create a tunnel and configure routing</li><li>Start the tunnel</li></ol>'
  },
  "post-go-interface": {
    title: "Go Composition and Method Overriding Pitfalls",
    date: "2024-09-04",
    category: "Golang",
    tags: ["Design Patterns"],
    readTime: "5 minute read",
    content: '<p>In a previous article, we discussed how to leverage composition in Go. This time, let\'s look at a common pitfall when trying to override methods in embedded structs.</p><h2>The Problem</h2><p>When struct A embeds struct B, and A wants to override one of B\'s methods, if B\'s method internally calls other methods, those calls will still use B\'s version, not A\'s.</p><h2>Solution</h2><p>Use interfaces to ensure method dispatch goes through the interface.</p>'
  },
  "post-black-wukong": {
    title: "My First Impressions of Black Myth: Wukong",
    date: "2024-08-21",
    category: "Game",
    tags: [],
    readTime: "8 minute read",
    content: '<p>August 20th, 2024, will forever be etched in the annals of Chinese gaming history. Black Myth: Wukong, developed by Game Science, finally launched after years of anticipation.</p><h2>Visuals</h2><p>The graphics are absolutely stunning. Unreal Engine 5 delivers incredible detail and lighting effects across every environment.</p><h2>Combat</h2><p>The combat feels tight and responsive. The transformations system is particularly creative and fun to experiment with.</p><h2>Overall</h2><p>As China\'s first AAA title, Black Myth: Wukong has delivered an impressive debut. A must-play for action RPG fans.</p>'
  },
  "post-hugo-theme": {
    title: "Quickly Build Blog with Hugo-theme-stack",
    date: "2024-02-20",
    category: "Web",
    tags: ["Hugo"],
    readTime: "5 minute read",
    content: '<p>Hugo is a static site generator known for its blazing fast build speeds. Combined with the Hugo Theme Stack, you can have a beautiful, modern blog up and running in minutes.</p><h2>Quick Start</h2><ol><li>Install Hugo</li><li>Create a new site</li><li>Add the Stack theme</li><li>Configure parameters</li><li>Write posts and deploy</li></ol>'
  },
  "post-jellyfin-02": {
    title: "Setting Up Jellyfin Media Server on NAS (Part 2)",
    date: "2023-09-16",
    category: "Knowledge Share",
    tags: ["NAS", "Jellyfin"],
    readTime: "6 minute read",
    content: '<p>After completing the installation and setup in the first part, let\'s look at some advanced configurations.</p><h2>Hardware Transcoding</h2><p>Configure Intel Quick Sync or NVIDIA NVENC for hardware-accelerated transcoding, significantly reducing CPU usage.</p><h2>Metadata Scraping</h2><p>Set up TheMovieDB, TMDB and other metadata sources for automatic movie and TV show information retrieval.</p>'
  },
  "post-jellyfin-01": {
    title: "Setting Up Jellyfin Media Server on NAS (Part 1)",
    date: "2023-09-15",
    category: "Knowledge Share",
    tags: ["NAS", "Jellyfin"],
    readTime: "5 minute read",
    content: '<p>Jellyfin is a free, open-source media server that lets you manage and stream your movies, TV shows, music, and more.</p><h2>Why Jellyfin?</h2><ul><li>Completely free and open-source</li><li>Multi-platform clients available</li><li>Rich feature set and plugin ecosystem</li></ul><h2>Docker Deployment</h2><pre><code>docker run -d --name=jellyfin -p 8096:8096 -v /path/to/config:/config -v /path/to/media:/media jellyfin/jellyfin</code></pre>'
  },
  "post-go-generic": {
    title: "Learn Go Generic",
    date: "2023-07-10",
    category: "Golang",
    tags: ["Go", "Generic"],
    readTime: "6 minute read",
    content: '<p>Go 1.18 introduced generics, a major feature update for the language.</p><h2>Basic Syntax</h2><pre><code>func Map[T any, U any](s []T, f func(T) U) []U {\n    result := make([]U, len(s))\n    for i, v := range s {\n        result[i] = f(v)\n    }\n    return result\n}</code></pre><h2>Use Cases</h2><ul><li>Generic data structures (stacks, queues, linked lists)</li><li>Utility functions (Map, Filter, Reduce)</li><li>Database operation wrappers</li></ul>'
  }
};

function getPostId(href) {
  return href.replace("#", "");
}

const root = document.documentElement;
const searchModal = document.querySelector(".search-modal");
const searchInput = document.querySelector("#searchInput");
const searchResults = document.querySelector("#searchResults");
const themeButtons = document.querySelectorAll(".theme-toggle");
const menuToggle = document.querySelector(".menu-toggle");

let themeMode = localStorage.getItem("blog-theme") || "light";

function applyTheme(mode) {
  themeMode = mode;
  localStorage.setItem("blog-theme", mode);
  root.dataset.scheme = mode;
  const label = mode === "dark" ? "亮色模式" : "暗色模式";
  const icon = mode === "dark" ? "☾" : "☀";
  themeButtons.forEach((button) => {
    if (button.classList.contains("language-button")) {
      button.innerHTML = '<span>◐</span>' + label;
    } else {
      button.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path></svg>' + label;
    }
    button.setAttribute("aria-label", label);
  });
}

function nextTheme() {
  applyTheme(themeMode === "dark" ? "light" : "dark");
}

function renderResults(query) {
  query = query || "";
  const normalized = query.trim().toLowerCase();
  const matches = articles.filter((article) => {
    const haystack = (article.title + " " + article.description).toLowerCase();
    return !normalized || haystack.includes(normalized);
  });
  searchResults.innerHTML = matches.map((article) => '<a class="search-result" href="' + article.href + '"><h3>' + article.title + '</h3><p>' + article.description + '</p></a>').join("") || '<div class="search-result"><h3>No results found</h3><p>Try a different keyword.</p></div>';
}

function openSearch(event) {
  if (event) event.preventDefault();
  searchModal.classList.add("open");
  searchModal.setAttribute("aria-hidden", "false");
  renderResults(searchInput.value);
  setTimeout(function() { searchInput.focus(); }, 40);
}

function closeSearch() {
  searchModal.classList.remove("open");
  searchModal.setAttribute("aria-hidden", "true");
}

var navLinks = document.querySelectorAll(".menu a");
var pageViews = document.querySelectorAll(".page-view");

var viewMap = {
  home: "view-home",
  archives: "view-archives",
  links: "view-links",
  media: "view-media",
  about: "view-about",
  post: "view-post"
};

function switchPage(id) {
  navLinks.forEach(function(link) {
    var href = link.getAttribute("href");
    link.classList.toggle("active", href === "#" + id);
  });
  var target = viewMap[id] || "view-home";
  pageViews.forEach(function(v) {
    v.style.display = v.id === target ? "" : "none";
  });
}

navLinks.forEach(function(link) {
  link.addEventListener("click", function(e) {
    e.preventDefault();
    var id = link.getAttribute("href").replace("#", "");
    switchPage(id);
    window.scrollTo({ top: 0, behavior: "smooth" });
    document.body.classList.remove("menu-open");
  });
});

document.querySelectorAll(".year-btn").forEach(function(btn) {
  btn.addEventListener("click", function() {
    document.querySelectorAll(".year-btn").forEach(function(b) { b.classList.remove("active"); });
    btn.classList.add("active");
  });
});

function showPost(href) {
  var id = getPostId(href);
  var post = postDB[id];
  if (!post) return;

  var detail = document.getElementById("postDetail");
  var tagsHtml = post.tags.map(function(t) { return '<span class="post-tag" onclick="searchByTag(\'' + t + '\')" style="cursor:pointer">' + t + '</span>'; }).join(" ");

  detail.innerHTML = '<h1>' + post.title + '</h1><div class="post-meta"><span><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"></rect><line x1="16" x2="16" y1="2" y2="6"></line><line x1="8" x2="8" y1="2" y2="6"></line><line x1="3" x2="21" y1="10" y2="10"></line></svg> ' + post.date + '</span><span><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg> ' + post.readTime + '</span><span>' + tagsHtml + '</span></div><div class="post-body">' + post.content + '</div>';

  pageViews.forEach(function(v) { v.style.display = "none"; });
  document.getElementById("view-post").style.display = "";
  navLinks.forEach(function(link) { link.classList.remove("active"); });
  window.scrollTo({ top: 0, behavior: "smooth" });

  // 生成目录
  generateToc();

  // 生成相关文章
  generateRelated(id);

  // 初始化 Waline 评论
  initWaline(id);
}

/* ===== 按标签搜索 ===== */
function searchByTag(tag) {
  // 打开搜索弹窗
  searchModal.classList.add("open");
  searchModal.setAttribute("aria-hidden", "false");
  searchInput.value = tag;
  searchInput.focus();

  // 执行搜索
  var results = Object.keys(postDB).filter(function(key) {
    return postDB[key].tags.some(function(t) {
      return t.toLowerCase().includes(tag.toLowerCase());
    });
  });

  searchResults.innerHTML = results.map(function(key) {
    var p = postDB[key];
    return '<a class="search-result" href="#' + key + '" onclick="event.preventDefault(); searchModal.classList.remove(\'open\'); showPost(\'#' + key + '\')"><h3>' + p.title + '</h3><p>' + p.description + '</p></a>';
  }).join("");

  if (results.length === 0) {
    searchResults.innerHTML = '<p style="color:var(--card-text-color-tertiary);text-align:center;padding:20px;">未找到包含标签 "' + tag + '" 的文章</p>';
  }
}

/* ===== 按分类搜索 ===== */
function searchByCategory(category) {
  // 打开搜索弹窗
  searchModal.classList.add("open");
  searchModal.setAttribute("aria-hidden", "false");
  searchInput.value = category;
  searchInput.focus();

  // 执行搜索
  var results = Object.keys(postDB).filter(function(key) {
    return postDB[key].category.toLowerCase().includes(category.toLowerCase());
  });

  searchResults.innerHTML = results.map(function(key) {
    var p = postDB[key];
    return '<a class="search-result" href="#' + key + '" onclick="event.preventDefault(); searchModal.classList.remove(\'open\'); showPost(\'#' + key + '\')"><h3>' + p.title + '</h3><p>' + p.description + '</p></a>';
  }).join("");

  if (results.length === 0) {
    searchResults.innerHTML = '<p style="color:var(--card-text-color-tertiary);text-align:center;padding:20px;">未找到分类为 "' + category + '" 的文章</p>';
  }
}

/* ===== 相关文章 ===== */
function generateRelated(currentId) {
  var relatedScroll = document.getElementById("relatedScroll");
  if (!relatedScroll) return;

  // 获取同分类的文章，排除当前文章
  var related = Object.keys(postDB)
    .filter(function(key) { return key !== currentId && postDB[key].category === postDB[currentId].category; })
    .slice(0, 5);

  // 如果同分类文章不足，补充其他文章
  if (related.length < 3) {
    Object.keys(postDB).forEach(function(key) {
      if (key !== currentId && related.indexOf(key) === -1) {
        related.push(key);
      }
      if (related.length >= 5) return;
    });
  }

  relatedScroll.innerHTML = related.map(function(key) {
    var p = postDB[key];
    return '<a class="related-card" href="#' + key + '" onclick="event.preventDefault(); showPost(\'#' + key + '\')"><h3>' + p.title + '</h3><div class="related-card-date">' + p.date + '</div></a>';
  }).join("");
}

/* ===== Waline 评论 ===== */
var walineLoaded = false;

function initWaline(postId) {
  var walineEl = document.getElementById("waline");
  if (!walineEl) return;

  // 清空旧内容
  walineEl.innerHTML = "";

  if (walineLoaded && window.Waline && typeof window.Waline.destroy === "function") {
    window.Waline.destroy();
  }

  var loadWaline = function() {
    if (walineLoaded && window.Waline) {
      renderWaline(postId);
      return;
    }

    // 加载 CSS
    var link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = "https://unpkg.com/@waline/client@v2/dist/waline.css";
    document.head.appendChild(link);

    // 加载 JS
    var script = document.createElement("script");
    script.src = "https://unpkg.com/@waline/client@v2/dist/waline.js";
    script.onload = function() {
      walineLoaded = true;
      renderWaline(postId);
    };
    document.body.appendChild(script);
  };

  loadWaline();
}

function renderWaline(postId) {
  var walineEl = document.getElementById("waline");
  if (!walineEl || !window.Waline) return;

  window.Waline.init({
    el: walineEl,
    serverURL: "https://comment.fcyan.xyz/",
    path: "/" + postId,
    lang: "zh-CN",
    dark: 'html[data-scheme="dark"]',
    meta: ["nick", "mail"],
    requiredMeta: ["nick"],
    pageSize: 10,
    reaction: true
  });
}

/* ===== 文章目录 (TOC) ===== */
var tocScrollHandler = null;

function generateToc() {
  var tocNav = document.getElementById("tocNav");
  var postBody = document.querySelector(".post-body");
  if (!tocNav || !postBody) return;

  // 移除旧的滚动监听
  if (tocScrollHandler) {
    window.removeEventListener("scroll", tocScrollHandler);
    tocScrollHandler = null;
  }

  // 收集所有 h2 和 h3 标题
  var headings = postBody.querySelectorAll("h2, h3");
  if (headings.length === 0) {
    document.getElementById("postToc").style.display = "none";
    return;
  }
  document.getElementById("postToc").style.display = "";

  // 为标题添加 id 并生成目录 HTML
  var ol = document.createElement("ol");
  var currentH2Li = null;
  var currentH2Ol = null;

  headings.forEach(function(heading, index) {
    // 生成标题 id
    var slug = "toc-" + index;
    heading.id = slug;

    if (heading.tagName === "H2") {
      currentH2Li = document.createElement("li");
      var a = document.createElement("a");
      a.href = "#" + slug;
      a.textContent = heading.textContent;
      a.dataset.target = slug;
      currentH2Li.appendChild(a);
      ol.appendChild(currentH2Li);
      currentH2Ol = null;
    } else if (heading.tagName === "H3") {
      if (!currentH2Li) {
        currentH2Li = document.createElement("li");
        ol.appendChild(currentH2Li);
      }
      if (!currentH2Ol) {
        currentH2Ol = document.createElement("ol");
        currentH2Li.appendChild(currentH2Ol);
      }
      var li = document.createElement("li");
      var a = document.createElement("a");
      a.href = "#" + slug;
      a.textContent = heading.textContent;
      a.dataset.target = slug;
      li.appendChild(a);
      currentH2Ol.appendChild(li);
    }
  });

  tocNav.innerHTML = "";
  tocNav.appendChild(ol);

  // 点击目录链接：平滑滚动
  tocNav.addEventListener("click", function(e) {
    var link = e.target.closest("a");
    if (!link) return;
    e.preventDefault();
    var target = document.getElementById(link.dataset.target);
    if (target) {
      target.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  });

  // 滚动高亮当前标题
  tocScrollHandler = function() {
    var scrollTop = window.scrollY + 120;
    var activeSlug = null;

    headings.forEach(function(heading) {
      if (heading.offsetTop <= scrollTop) {
        activeSlug = heading.id;
      }
    });

    tocNav.querySelectorAll("a").forEach(function(a) {
      a.classList.toggle("active", a.dataset.target === activeSlug);
    });
  };
  window.addEventListener("scroll", tocScrollHandler);
  tocScrollHandler();
}

document.addEventListener("click", function(e) {
  var link = e.target.closest('a[href^="#post-"]');
  if (link) {
    e.preventDefault();
    var href = link.getAttribute("href");
    if (postDB[getPostId(href)]) {
      showPost(href);
    }
  }
});

document.getElementById("backToList").addEventListener("click", function() {
  switchPage("home");
});

themeButtons.forEach(function(button) { button.addEventListener("click", nextTheme); });
document.querySelectorAll("[data-open-search]").forEach(function(button) { button.addEventListener("click", openSearch); });
document.querySelectorAll("[data-close-search]").forEach(function(button) { button.addEventListener("click", closeSearch); });
searchInput.addEventListener("input", function() { renderResults(searchInput.value); });
menuToggle.addEventListener("click", function() { document.body.classList.toggle("menu-open"); });
document.addEventListener("keydown", function(event) {
  if (event.key === "Escape") {
    closeSearch();
    document.body.classList.remove("menu-open");
  }
});

applyTheme(themeMode);
renderResults();
switchPage("home");

// ===== 回到顶部按钮（带进度） =====
const backToTopBtn = document.getElementById('backToTop');
const progressBar = document.getElementById('progressBar');
const circumference = 2 * Math.PI * 18; // r=18

function updateProgress() {
  const scrollTop = window.scrollY;
  const docHeight = document.documentElement.scrollHeight - window.innerHeight;
  const progress = docHeight > 0 ? Math.min(scrollTop / docHeight, 1) : 0;

  // 更新进度条
  const offset = circumference - (progress * circumference);
  progressBar.style.strokeDashoffset = offset;

  // 只在文章详情页面显示回到顶部按钮
  const isPostView = document.getElementById('view-post') && document.getElementById('view-post').style.display !== 'none';

  // 显示/隐藏按钮
  if (scrollTop > 300 && isPostView) {
    backToTopBtn.classList.add('visible');
  } else {
    backToTopBtn.classList.remove('visible');
  }
}

window.addEventListener('scroll', updateProgress);
updateProgress();

backToTopBtn.addEventListener('click', function() {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
});

// ===== 代码块自动折叠 =====
function initCodeCollapse() {
  document.querySelectorAll('.post-body pre').forEach(function(pre) {
    // 避免重复处理
    if (pre.parentElement.classList.contains('code-collapse-wrapper')) return;

    // 只对超过 600px 的代码块折叠
    if (pre.scrollHeight <= 600) return;

    var wrapper = document.createElement('div');
    wrapper.className = 'code-collapse-wrapper';
    pre.parentNode.insertBefore(wrapper, pre);
    wrapper.appendChild(pre);

    var btn = document.createElement('button');
    btn.className = 'code-collapse-btn';
    btn.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg> 展开代码';
    wrapper.parentNode.insertBefore(btn, wrapper.nextSibling);

    btn.addEventListener('click', function() {
      var isExpanded = wrapper.classList.contains('expanded');
      if (isExpanded) {
        wrapper.classList.remove('expanded');
        btn.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg> 展开代码';
        // 收起时滚动回代码块顶部
        wrapper.scrollIntoView({ behavior: 'smooth', block: 'start' });
      } else {
        wrapper.classList.add('expanded');
        btn.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="18 15 12 9 6 15"></polyline></svg> 收起代码';
      }
    });
  });
}

// 在文章切换时重新初始化代码折叠
var originalShowPost = showPost;
showPost = function(href) {
  originalShowPost(href);
  setTimeout(initCodeCollapse, 50);
};

// ===== 博客运行时长 =====
function initSiteRuntime() {
  var runtimeEl = document.getElementById('siteRuntime');
  if (!runtimeEl) return;

  // 建站日期（可根据实际情况修改）
  var siteStartDate = new Date('2024-02-20');

  function updateRuntime() {
    var now = new Date();
    var diff = now - siteStartDate;

    var days = Math.floor(diff / (1000 * 60 * 60 * 24));
    var hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

    runtimeEl.innerHTML = '<span class="site-runtime-icon">🕐</span> 已运行 <span class="site-runtime-number">' + days + '</span> 天 <span class="site-runtime-number">' + hours + '</span> 小时 <span class="site-runtime-number">' + minutes + '</span> 分钟';
  }

  updateRuntime();
  // 每分钟更新一次
  setInterval(updateRuntime, 60000);
}

// 初始化
initCodeCollapse();
initSiteRuntime();
