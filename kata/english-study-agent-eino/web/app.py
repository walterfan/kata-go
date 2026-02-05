import streamlit as st
import streamlit.components.v1 as components
import requests
import re
from datetime import datetime

API_URL = "http://localhost:8080/api"

def speak_text(text: str, rate: float = 0.9):
    """Create an HTML component that speaks text using Web Speech API."""
    # Escape the text for JavaScript
    escaped_text = text.replace("\\", "\\\\").replace("'", "\\'").replace('"', '\\"').replace("\n", " ").replace("\r", "")
    
    html_code = f'''
    <div id="tts-container" style="display: none;">
        <script>
            // Cancel any ongoing speech first
            window.speechSynthesis.cancel();
            
            // Create utterance
            var utterance = new SpeechSynthesisUtterance("{escaped_text}");
            utterance.lang = 'en-US';
            utterance.rate = {rate};
            
            // Try to get a good English voice
            var voices = window.speechSynthesis.getVoices();
            if (voices.length > 0) {{
                var englishVoice = voices.find(v => v.lang.startsWith('en'));
                if (englishVoice) {{
                    utterance.voice = englishVoice;
                }}
            }}
            
            // Speak
            window.speechSynthesis.speak(utterance);
        </script>
    </div>
    '''
    components.html(html_code, height=0)

def stop_speech():
    """Create an HTML component that stops speech."""
    html_code = '''
    <div id="tts-stop" style="display: none;">
        <script>
            window.speechSynthesis.cancel();
        </script>
    </div>
    '''
    components.html(html_code, height=0)

st.set_page_config(
    page_title="Language Learning | 语言学习", 
    page_icon="🌍", 
    layout="wide",
    initial_sidebar_state="expanded"
)

# Internationalization (i18n)
LANG = {
    "en": {
        "title": "📚 English Agent",
        "subtitle": "AI-powered English learning",
        "main_title": "The limits of my language mean the limits of my world.",
        "main_subtitle": "— Ludwig Wittgenstein | AI-powered English learning",
        "input_mode": "📂 Input Mode",
        "choose_source": "Choose input source:",
        "article_mode": "📰 Article",
        "text_mode": "✏️ Text Input",
        "news_articles": "📰 News Articles",
        "refresh": "🔄 Refresh Headlines",
        "fetching": "Fetching...",
        "loaded": "loaded!",
        "no_articles_found": "⚠️ No articles found. The RSS feed might be empty or unavailable.",
        "select_article": "Select an article:",
        "preview": "📄 Preview",
        "click_refresh": "👆 Click Refresh to load articles",
        "text_input": "✏️ Text Input",
        "enter_text": "Enter text:",
        "paste_here": "Paste English text here...",
        "words": "words",
        "examples": "💡 Examples",
        "actions": "🎯 Actions",
        "select_first": "Select an article or enter text first",
        "explain": "📖 Explain",
        "summarize": "📝 Summarize",
        "translate": "🌐 Translate",
        "refine": "✨ Refine",
        "sentences": "📋 Sentences",
        "vocabulary": "📚 Vocabulary",
        "clear": "🗑️ Clear Results",
        "current_text": "📄 Current Text",
        "use_sidebar": "👈 Use the sidebar to load an article or enter text, then click an action button.",
        "results": "💬 Results",
        "welcome": "👋 Welcome!",
        "welcome_msg": "Ready to learn English with AI assistance.",
        "quick_start": "Quick Start:",
        "step1": "1️⃣ Choose input mode in the sidebar (Article or Text)",
        "step2": "2️⃣ Load content or paste your text",
        "step3": "3️⃣ Click an action button to analyze",
        "recommended": "💡 Recommended Flow:",
        "flow": "Explain → Vocabulary → Sentences",
        "tip": "Tip: Click the ◀ button at top-left to hide/show the sidebar",
        "footer1": "📚 English Learning Agent v1.0",
        "footer2": "Built with CloudWeGo Eino + Streamlit",
        "footer3": "🇨🇳 Supports Chinese Translation",
        "language": "🌐 Language",
        "backend_error": "❌ Backend not running",
        "processing": "Processing...",
        "no_text_warning": "⚠️ Please enter some text or select an article first!",
        "error": "Error:",
        "timeout": "⏱️ Request timed out. Please try again.",
        "connect_error": "❌ Cannot connect to backend. Please run `make run` first.",
        # Help tooltips
        "help_explain": "Sentence by sentence",
        "help_summarize": "Concise summary",
        "help_translate": "To Chinese",
        "help_refine": "Improve text",
        "help_sentences": "Extract patterns",
        "help_vocabulary": "Key words",
        # RSS source
        "rss_source": "📡 RSS Source",
        "select_source": "Select source:",
        "all_sources": "All Sources",
        "load_sources": "Loading sources...",
        "no_sources": "No RSS sources configured",
        # Settings tab
        "settings": "⚙️ Settings",
        "settings_mode": "⚙️ Settings",
        "manage_feeds": "Manage RSS Feeds",
        "default_feeds": "📋 Default Feeds (from config)",
        "custom_feeds": "✏️ Custom Feeds",
        "add_feed": "➕ Add New Feed",
        "feed_title": "Title",
        "feed_url": "URL",
        "feed_category": "Category",
        "feed_enabled": "Enabled",
        "save_feed": "💾 Save",
        "delete_feed": "🗑️ Delete",
        "edit_feed": "✏️ Edit",
        "cancel": "Cancel",
        "feed_added": "✅ Feed added successfully!",
        "feed_updated": "✅ Feed updated successfully!",
        "feed_deleted": "✅ Feed deleted successfully!",
        "no_custom_feeds": "No custom feeds added yet",
        "category_english": "English Learning",
        "category_tech": "Technology",
        "category_medical": "Medical",
        "category_news": "News",
        "category_other": "Other",
        # TTS (Text-to-Speech)
        "read_article": "🔊 Read",
        "stop_reading": "⏹ Stop",
        "help_read": "Read the article aloud",
        "help_stop": "Stop reading",
        "reading": "Reading...",
        # URL Fetch
        "fetch_url": "🔗 Fetch from URL",
        "enter_url": "Enter article URL:",
        "url_placeholder": "https://example.com/article...",
        "load_url": "📥 Load Article",
        "loading_url": "Loading article...",
        "url_loaded": "Article loaded!",
        "url_error": "Failed to fetch article",
        "or_divider": "— OR —",
        # Streaming
        "enable_streaming": "⚡ Streaming Mode",
        "streaming_hint": "See AI response in real-time",
    },
    "zh": {
        "title": "📚 英语学习助手",
        "subtitle": "AI驱动的英语学习工具",
        "main_title": "我语言的边界，就是我世界的边界。",
        "main_subtitle": "— 路德维希·维特根斯坦 | AI驱动的英语学习",
        "input_mode": "📂 输入模式",
        "choose_source": "选择输入来源：",
        "article_mode": "📰 文章",
        "text_mode": "✏️ 文本输入",
        "news_articles": "📰 新闻文章",
        "refresh": "🔄 刷新文章",
        "fetching": "获取中...",
        "loaded": "已加载！",
        "no_articles_found": "⚠️ 未找到文章。RSS 源可能为空或不可用。",
        "select_article": "选择文章：",
        "preview": "📄 预览",
        "click_refresh": "👆 点击刷新加载文章",
        "text_input": "✏️ 文本输入",
        "enter_text": "输入文本：",
        "paste_here": "在此粘贴英文文本...",
        "words": "词",
        "examples": "💡 示例",
        "actions": "🎯 操作",
        "select_first": "请先选择文章或输入文本",
        "explain": "📖 逐句解释",
        "summarize": "📝 摘要",
        "translate": "🌐 翻译",
        "refine": "✨ 润色",
        "sentences": "📋 句型",
        "vocabulary": "📚 词汇",
        "clear": "🗑️ 清除结果",
        "current_text": "📄 当前文本",
        "use_sidebar": "👈 使用侧边栏加载文章或输入文本，然后点击操作按钮。",
        "results": "💬 结果",
        "welcome": "👋 欢迎！",
        "welcome_msg": "准备好使用AI助手学习英语了。",
        "quick_start": "快速开始：",
        "step1": "1️⃣ 在侧边栏选择输入模式（文章或文本）",
        "step2": "2️⃣ 加载内容或粘贴文本",
        "step3": "3️⃣ 点击操作按钮进行分析",
        "recommended": "💡 推荐流程：",
        "flow": "逐句解释 → 词汇 → 句型",
        "tip": "提示：点击左上角的 ◀ 按钮可以隐藏/显示侧边栏",
        "footer1": "📚 英语学习助手 v1.0",
        "footer2": "基于 CloudWeGo Eino + Streamlit 构建",
        "footer3": "🇨🇳 支持中文翻译",
        "language": "🌐 语言",
        "backend_error": "❌ 后端未运行",
        "processing": "处理中...",
        "no_text_warning": "⚠️ 请先输入文本或选择文章！",
        "error": "错误：",
        "timeout": "⏱️ 请求超时，请重试。",
        "connect_error": "❌ 无法连接后端，请先运行 `make run`。",
        # Help tooltips
        "help_explain": "逐句解释含义",
        "help_summarize": "简洁摘要",
        "help_translate": "翻译成中文",
        # RSS source
        "rss_source": "📡 RSS 源",
        "select_source": "选择来源：",
        "all_sources": "所有来源",
        "load_sources": "加载来源中...",
        "no_sources": "未配置RSS源",
        "help_refine": "改进文本",
        "help_sentences": "提取句型",
        "help_vocabulary": "关键词汇",
        # Settings tab
        "settings": "⚙️ 设置",
        "settings_mode": "⚙️ 设置",
        "manage_feeds": "管理RSS订阅",
        "default_feeds": "📋 默认订阅（来自配置文件）",
        "custom_feeds": "✏️ 自定义订阅",
        "add_feed": "➕ 添加新订阅",
        "feed_title": "标题",
        "feed_url": "URL地址",
        "feed_category": "分类",
        "feed_enabled": "启用",
        "save_feed": "💾 保存",
        "delete_feed": "🗑️ 删除",
        "edit_feed": "✏️ 编辑",
        "cancel": "取消",
        "feed_added": "✅ 订阅添加成功！",
        "feed_updated": "✅ 订阅更新成功！",
        "feed_deleted": "✅ 订阅删除成功！",
        "no_custom_feeds": "尚未添加自定义订阅",
        "category_english": "英语学习",
        "category_tech": "科技",
        "category_medical": "医学",
        "category_news": "新闻",
        "category_other": "其他",
        # TTS (Text-to-Speech)
        "read_article": "🔊 朗读",
        "stop_reading": "⏹ 停止",
        "help_read": "朗读文章",
        "help_stop": "停止朗读",
        "reading": "朗读中...",
        # URL Fetch
        "fetch_url": "🔗 从URL加载",
        "enter_url": "输入文章链接：",
        "url_placeholder": "https://example.com/article...",
        "load_url": "📥 加载文章",
        "loading_url": "加载文章中...",
        "url_loaded": "文章已加载！",
        "url_error": "获取文章失败",
        "or_divider": "— 或者 —",
        # Streaming
        "enable_streaming": "⚡ 流式响应",
        "streaming_hint": "实时查看AI生成内容",
    }
}

def t(key: str) -> str:
    """Get translated text for the current language."""
    lang = st.session_state.get("language", "en")
    return LANG.get(lang, LANG["en"]).get(key, key)

def extract_links(text: str) -> list[tuple[str, str]]:
    """Extract links from text (Markdown and HTML). Returns list of (text, url)."""
    links = []
    
    # Extract Markdown links: [text](url)
    md_links = re.findall(r'\[([^\]]+)\]\((https?://[^)]+)\)', text)
    links.extend(md_links)
    
    # Extract HTML links: <a href="url">text</a>
    # Note: simple regex, might not handle all attributes perfectly
    html_links = re.findall(r'<a\s+(?:[^>]*?\s+)?href="([^"]*)"[^>]*>(.*?)</a>', text, re.IGNORECASE)
    # Regex returns (url, text), swap to (text, url)
    links.extend([(t, u) for u, t in html_links])
    
    # Deduplicate by URL
    seen = set()
    unique_links = []
    for text, url in links:
        if url not in seen:
            seen.add(url)
            # clean text: remove html tags if any inside link text
            text = re.sub(r'<[^>]+>', '', text).strip()
            if not text:
                text = url
            # Truncate long text
            if len(text) > 50:
                text = text[:47] + "..."
            unique_links.append((text, url))
            
    return unique_links

# ... existing code ...

# Custom CSS for better styling
st.markdown("""
<style>
    /* Main content area */
    .block-container {
        padding-top: 1rem;
        padding-bottom: 1rem;
        max-width: 100%;
    }
    
    /* Sidebar styling */
    [data-testid="stSidebar"] {
        min-width: 320px;
        max-width: 400px;
    }
    
    [data-testid="stSidebar"] .block-container {
        padding-top: 1rem;
    }
    
    /* Tab styling */
    .stTabs [data-baseweb="tab-list"] {
        gap: 8px;
    }
    .stTabs [data-baseweb="tab"] {
        padding: 10px 20px;
        font-size: 16px;
    }
    
    /* Article card */
    .article-card {
        background-color: #f8f9fa;
        padding: 15px;
        border-radius: 10px;
        border: 1px solid #e0e0e0;
        margin: 10px 0;
    }
    
    /* Chat message styling */
    [data-testid="stChatMessage"] {
        max-width: 100%;
    }
    
    /* English sentences in blockquotes - larger font */
    [data-testid="stChatMessage"] blockquote {
        font-size: 1.2em;
        line-height: 1.6;
        border-left: 4px solid #667eea;
        background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 100%);
        padding: 12px 16px;
        margin: 10px 0;
        border-radius: 0 8px 8px 0;
    }
    
    /* Chinese translation styling */
    [data-testid="stChatMessage"] p:has(> span:first-child) {
        font-size: 1.05em;
        color: #444;
    }
    
    /* Text area styling */
    .stTextArea textarea {
        font-size: 16px;
    }
    
    /* Button row spacing */
    .stButton button {
        margin-bottom: 5px;
    }
    
    /* Compact action buttons - no text wrap */
    [data-testid="stHorizontalBlock"] .stButton button {
        padding: 0.4rem 0.8rem;
        font-size: 0.9rem;
        white-space: nowrap;
        min-height: 0;
    }
    
    /* Reduce column gap for action buttons */
    [data-testid="stHorizontalBlock"] {
        gap: 0.5rem;
    }
    
    /* Title for proverb - no wrap */
    h1 {
        font-size: 1.8rem !important;
        line-height: 1.6 !important;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }
    
    /* Action button section */
    .action-section {
        background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
        padding: 15px;
        border-radius: 10px;
        margin: 10px 0;
    }
    
</style>
""", unsafe_allow_html=True)

# Task definitions with Chinese support
TASKS = {
    "explain": {
        "icon": "📖",
        "label": "Explain",
        "label_zh": "逐句解释",
        "description": "Explain sentence by sentence",
        "prompt": """Explain this text sentence by sentence. For each sentence, use this EXACT format:

---

### 📝 Sentence 1

> **[English sentence here - show the original sentence]**

🇨🇳 [Chinese translation here - 中文翻译]

💡 **Key phrases:** [list important words/phrases with brief explanations]

---

### 📝 Sentence 2

> **[English sentence here]**

🇨🇳 [Chinese translation]

💡 **Key phrases:** [important words/phrases]

---

(Continue this format for ALL sentences in the text)

IMPORTANT formatting rules:
1. Put English sentence in a blockquote with bold text (use > ** **)
2. Put Chinese translation on a NEW line starting with 🇨🇳
3. Put key phrases on a NEW line starting with 💡
4. Use --- between each sentence block
5. Number each sentence (Sentence 1, Sentence 2, etc.)"""
    },
    "summarize": {
        "icon": "📝",
        "label": "Summarize",
        "label_zh": "摘要",
        "description": "Get a concise summary",
        "prompt": """Summarize this text:
1. Provide a brief summary in English (3-5 sentences)
2. List the main points
3. Provide Chinese summary (中文摘要)

Format:
**English Summary:**
[summary]

**Main Points:**
- [point 1]
- [point 2]
- [point 3]

**中文摘要:**
[Chinese summary]"""
    },
    "translate": {
        "icon": "🌐",
        "label": "Translate",
        "label_zh": "翻译",
        "description": "Translate to Chinese",
        "prompt": """Translate this text to Chinese:
1. Provide accurate Chinese translation
2. Keep the original tone and style
3. Add notes for any culturally-specific terms

Format:

**Original English:**
[original text]

**中文翻译:**
[Chinese translation]

**Translation Notes:**
[any important notes about the translation]"""
    },
    "refine": {
        "icon": "✨",
        "label": "Refine",
        "label_zh": "润色",
        "description": "Improve and simplify",
        "prompt": """Refine and improve this text:
1. Rewrite in clearer, simpler English
2. Fix any grammar issues
3. Suggest better word choices
4. Provide Chinese explanation of changes (中文说明修改内容)

Format:
**Original:**
[original text]

**Refined Version:**
[improved text]

**Changes Made:**
- [change 1 and why]
- [change 2 and why]

**中文说明:**
[Chinese explanation of what was changed and why]"""
    },
    "extract_sentences": {
        "icon": "📋",
        "label": "Sentences",
        "label_zh": "句型提取",
        "description": "Extract useful sentence patterns",
        "prompt": """Extract useful sentence structures from this text:
1. Identify 3-5 useful sentence patterns
2. Explain when to use each pattern
3. Provide example sentences
4. Include Chinese explanation (中文说明)

Format for each pattern:
---
**Pattern:** [sentence pattern with blanks]
**Example from text:** [original sentence]
**When to use:** [explanation]
**Your own example:** [create a new example]

**中文说明:** [Chinese explanation]
---"""
    },
    "extract_vocabulary": {
        "icon": "📚",
        "label": "Vocabulary",
        "label_zh": "词汇提取",
        "description": "Extract key words and phrases",
        "prompt": """Extract important vocabulary from this text:
1. List 5-8 key words and phrases
2. Provide definition, pronunciation hint, and example
3. Include Chinese translation (中文翻译)
4. Rate difficulty (Basic/Intermediate/Advanced)

Format for each word:
---
**Word/Phrase:** [word or phrase]
**Pronunciation:** [how to pronounce]
**Meaning:** [definition in simple English]
**中文:** [Chinese translation]
**Example:** [example sentence]
**Level:** [Basic/Intermediate/Advanced]
---"""
    }
}

# Initialize session state
if "messages" not in st.session_state:
    st.session_state.messages = []
if "current_text" not in st.session_state:
    st.session_state.current_text = ""
if "articles" not in st.session_state:
    st.session_state.articles = []
if "selected_article" not in st.session_state:
    st.session_state.selected_article = None
if "input_mode" not in st.session_state:
    st.session_state.input_mode = "article"
# Streaming state
if "enable_streaming" not in st.session_state:
    st.session_state.enable_streaming = True
if "streaming_task" not in st.session_state:
    st.session_state.streaming_task = None
if "streaming_text" not in st.session_state:
    st.session_state.streaming_text = None
if "streaming_timestamp" not in st.session_state:
    st.session_state.streaming_timestamp = None
if "language" not in st.session_state:
    st.session_state.language = "en"
if "rss_sources" not in st.session_state:
    st.session_state.rss_sources = []
if "selected_source" not in st.session_state:
    st.session_state.selected_source = "all"

def call_agent(text: str, task_key: str) -> str:
    """Call the backend agent with the given text and task."""
    task = TASKS[task_key]
    try:
        response = requests.post(
            f"{API_URL}/chat",
            json={"text": text, "task": task["prompt"]},
            timeout=120
        )
        if response.status_code == 200:
            return response.json().get("result", "")
        else:
            return f"❌ {t('error')} {response.status_code}"
    except requests.exceptions.ConnectionError:
        return t("connect_error")
    except requests.exceptions.Timeout:
        return t("timeout")
    except Exception as e:
        return f"❌ {t('error')} {e}"

def call_agent_stream(text: str, task_key: str):
    """Call the backend agent with streaming response using SSE."""
    task = TASKS[task_key]
    try:
        response = requests.post(
            f"{API_URL}/chat/stream",
            json={"text": text, "task": task["prompt"]},
            stream=True,
            timeout=120
        )
        if response.status_code == 200:
            for line in response.iter_lines():
                if line:
                    line = line.decode('utf-8')
                    # Parse SSE format: "event: message\ndata: content"
                    if line.startswith('data:'):
                        # Remove "data:" prefix
                        data = line[5:]
                        # Remove optional leading space defined by SSE spec
                        if data.startswith(' '):
                            data = data[1:]
                        
                        # Don't strip() here! It removes leading/trailing spaces from content
                        
                        if data:
                            # Restore escaped newlines
                            data = data.replace('\\n', '\n')
                            yield data
                    elif line.startswith('event:'):
                        event = line[6:].strip()
                        if event == 'done':
                            break
                        elif event == 'error':
                            yield f"\n❌ Stream error"
                            break
        else:
            yield f"❌ {t('error')} {response.status_code}"
    except requests.exceptions.ConnectionError:
        yield t("connect_error")
    except requests.exceptions.Timeout:
        yield t("timeout")
    except Exception as e:
        yield f"❌ {t('error')} {e}"

def process_task(task_key: str, text: str):
    """Process a task with the given text."""
    if not text.strip():
        st.warning(t("no_text_warning"))
        return
    
    task = TASKS[task_key]
    lang = st.session_state.get("language", "en")
    label = task["label_zh"] if lang == "zh" else task["label"]
    
    # Get current timestamp
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    
    # Add user message first
    user_msg = f"{task['icon']} **{label}**"
    st.session_state.messages.append({
        "role": "user", 
        "content": user_msg,
        "timestamp": timestamp
    })
    
    # Check if streaming is enabled
    use_streaming = st.session_state.get("enable_streaming", True)
    
    if use_streaming:
        # Streaming mode - show response as it generates
        st.session_state.streaming_task = task_key
        st.session_state.streaming_text = text
        st.session_state.streaming_timestamp = timestamp
        st.rerun()
    else:
        # Non-streaming mode - wait for complete response
        with st.spinner(f"🤔 {t('processing')}"):
            result = call_agent(text, task_key)
        
        st.session_state.messages.append({
            "role": "assistant", 
            "content": result,
            "timestamp": timestamp,
            "task": task_key
        })
        st.rerun()

def render_content_with_tts(content: str, task_key: str, msg_index: int):
    """Render content with TTS buttons for English sentences."""
    # For explain task, extract sentences and add TTS buttons
    if task_key == "explain":
        # Extract English sentences from blockquotes
        sentences = re.findall(r'>\s*\*\*(.+?)\*\*', content)
        
        # First render the content normally
        st.markdown(content)
        
        # Then add TTS buttons for each sentence in an expander
        if sentences:
            with st.expander("🔊 Read sentences", expanded=False):
                for idx, sentence in enumerate(sentences):
                    col1, col2 = st.columns([6, 1])
                    with col1:
                        st.caption(f"{idx+1}. {sentence[:60]}..." if len(sentence) > 60 else f"{idx+1}. {sentence}")
                    with col2:
                        if st.button("🔊", key=f"tts_{msg_index}_{idx}", help=f"Read: {sentence[:30]}..."):
                            st.session_state.tts_text = sentence
                            st.session_state.tts_action = "speak"
    elif task_key == "extract_sentences":
        # Render sentence patterns with styled Streamlit components
        render_sentence_patterns(content, msg_index)
    elif task_key == "extract_vocabulary":
        # Render vocabulary with styled Streamlit components
        render_vocabulary(content, msg_index)
    else:
        # For other tasks, just render normally
        st.markdown(content)

def render_sentence_patterns(content: str, msg_index: int):
    """Render sentence patterns with styled Streamlit components."""
    # Split content by --- separator
    patterns = re.split(r'\n-{3,}\n', content)
    
    for idx, pattern in enumerate(patterns):
        if not pattern.strip():
            continue
        
        # Parse the pattern content
        pattern_match = re.search(r'\*\*Pattern:\*\*\s*(.+?)(?=\n\*\*|\Z)', pattern, re.DOTALL)
        example_match = re.search(r'\*\*Example from text:\*\*\s*(.+?)(?=\n\*\*|\Z)', pattern, re.DOTALL)
        when_match = re.search(r'\*\*When to use:\*\*\s*(.+?)(?=\n\*\*|\Z)', pattern, re.DOTALL)
        own_example_match = re.search(r'\*\*Your own example:\*\*\s*(.+?)(?=\n\*\*|\Z)', pattern, re.DOTALL)
        chinese_match = re.search(r'\*\*中文说明:\*\*\s*(.+?)(?=\n\*\*|\Z)', pattern, re.DOTALL)
        
        if pattern_match:
            pattern_text = pattern_match.group(1).strip()
            example_text = example_match.group(1).strip() if example_match else ""
            when_text = when_match.group(1).strip() if when_match else ""
            own_example_text = own_example_match.group(1).strip() if own_example_match else ""
            chinese_text = chinese_match.group(1).strip() if chinese_match else ""
            
            # Use Streamlit native components for better rendering
            with st.container():
                st.markdown(f"### 📋 句型 #{idx + 1}")
                st.info(f"**{pattern_text}**")
                
                col1, col2 = st.columns(2)
                with col1:
                    st.markdown("**📝 原文例句:**")
                    st.caption(f"*\"{example_text}\"*")
                    
                    st.markdown("**✍️ 自造例句:**")
                    st.caption(f"*\"{own_example_text}\"*")
                
                with col2:
                    st.markdown("**💡 使用场景:**")
                    st.caption(when_text)
                    
                    st.markdown("**🇨🇳 中文说明:**")
                    st.caption(chinese_text)
                
                st.divider()
        else:
            # Fallback to regular markdown if parsing fails
            st.markdown(pattern)

def render_vocabulary(content: str, msg_index: int):
    """Render vocabulary with styled Streamlit components."""
    # Split content by --- separator
    words = re.split(r'\n-{3,}\n', content)
    
    for idx, word_block in enumerate(words):
        if not word_block.strip():
            continue
        
        # Parse the word content
        word_match = re.search(r'\*\*Word/Phrase:\*\*\s*(.+?)(?=\n\*\*|\Z)', word_block, re.DOTALL)
        pron_match = re.search(r'\*\*Pronunciation:\*\*\s*(.+?)(?=\n\*\*|\Z)', word_block, re.DOTALL)
        meaning_match = re.search(r'\*\*Meaning:\*\*\s*(.+?)(?=\n\*\*|\Z)', word_block, re.DOTALL)
        chinese_match = re.search(r'\*\*中文:\*\*\s*(.+?)(?=\n\*\*|\Z)', word_block, re.DOTALL)
        example_match = re.search(r'\*\*Example:\*\*\s*(.+?)(?=\n\*\*|\Z)', word_block, re.DOTALL)
        level_match = re.search(r'\*\*Level:\*\*\s*(.+?)(?=\n\*\*|\Z)', word_block, re.DOTALL)
        
        if word_match:
            word_text = word_match.group(1).strip()
            pron_text = pron_match.group(1).strip() if pron_match else ""
            meaning_text = meaning_match.group(1).strip() if meaning_match else ""
            chinese_text = chinese_match.group(1).strip() if chinese_match else ""
            example_text = example_match.group(1).strip() if example_match else ""
            level_text = level_match.group(1).strip() if level_match else "Intermediate"
            
            # Determine level indicator
            level_indicators = {
                "Basic": "🟢 Basic",
                "Intermediate": "🟡 Intermediate", 
                "Advanced": "🔴 Advanced"
            }
            level_display = level_indicators.get(level_text, "🟡 Intermediate")
            
            # Use Streamlit native components
            with st.container():
                col_word, col_level = st.columns([4, 1])
                with col_word:
                    st.markdown(f"### 📚 {word_text}")
                    if pron_text:
                        st.caption(f"/{pron_text}/")
                with col_level:
                    st.markdown(f"**{level_display}**")
                
                st.markdown(f"**📖 Meaning:** {meaning_text}")
                st.markdown(f"**🇨🇳 中文:** {chinese_text}")
                st.markdown(f"**💬 Example:** *\"{example_text}\"*")
                
                st.divider()
        else:
            # Fallback to regular markdown if parsing fails
            st.markdown(word_block)

# ==================== SIDEBAR ====================
with st.sidebar:
    st.title(t("title"))
    st.caption(t("subtitle"))
    
    st.divider()
    
    # Language selector
    st.markdown(f"### {t('language')}")
    lang_options = {"English": "en", "中文": "zh"}
    current_lang = "中文" if st.session_state.language == "zh" else "English"
    selected_lang = st.radio(
        "Language:",
        list(lang_options.keys()),
        index=list(lang_options.keys()).index(current_lang),
        label_visibility="collapsed",
        key="lang_radio",
        horizontal=True
    )
    if lang_options[selected_lang] != st.session_state.language:
        st.session_state.language = lang_options[selected_lang]
        st.rerun()
    
    # Streaming toggle
    st.session_state.enable_streaming = st.toggle(
        t("enable_streaming"),
        value=st.session_state.enable_streaming,
        help=t("streaming_hint")
    )
    
    st.divider()
    
    # Mode selector
    st.markdown(f"### {t('input_mode')}")
    mode_options = [t("article_mode"), t("text_mode"), t("settings_mode")]
    mode = st.radio(
        t("choose_source"),
        mode_options,
        label_visibility="collapsed",
        key="mode_radio"
    )
    if mode == mode_options[0]:
        st.session_state.input_mode = "article"
    elif mode == mode_options[1]:
        st.session_state.input_mode = "text"
    else:
        st.session_state.input_mode = "settings"
    
    st.divider()
    
    # ========== ARTICLE MODE ==========
    if st.session_state.input_mode == "article":
        st.markdown(f"### {t('news_articles')}")
        
        # RSS Source selector
        st.markdown(f"**{t('rss_source')}**")

        # Load RSS sources if not loaded
        if not st.session_state.rss_sources:
            try:
                response = requests.get(f"{API_URL}/rss-sources", timeout=5)
                if response.status_code == 200:
                    st.session_state.rss_sources = response.json().get("sources", [])
            except:
                pass
        
        # Source selector dropdown
        if st.session_state.rss_sources:
            source_options = [t("all_sources")] + [s["title"] for s in st.session_state.rss_sources]
            selected_source_idx = st.selectbox(
                t("select_source"),
                range(len(source_options)),
                format_func=lambda x: source_options[x],
                key="rss_source_selector",
                label_visibility="collapsed"
            )
            st.session_state.selected_source = "all" if selected_source_idx == 0 else source_options[selected_source_idx]
        else:
            st.caption(t("no_sources"))
        
        # Refresh button
        if st.button(t("refresh"), use_container_width=True):
            with st.spinner(t("fetching")):
                try:
                    # Build URL with source parameter
                    url = f"{API_URL}/feeds"
                    if st.session_state.selected_source != "all":
                        url += f"?source={st.session_state.selected_source}"
                    
                    response = requests.get(url, timeout=15)
                    if response.status_code == 200:
                        articles = response.json().get("articles") or []
                        st.session_state.articles = articles
                        if articles:
                            st.success(f"✅ {len(articles)} {t('loaded')}")
                        else:
                            st.warning(t("no_articles_found"))
                except requests.exceptions.ConnectionError:
                    st.error(t("backend_error"))
                except Exception as e:
                    st.error(f"{t('error')} {e}")
        
        # Article list
        if st.session_state.articles:
            article_titles = [f"{a['Title'][:40]}..." for a in st.session_state.articles[:12]]
            selected_idx = st.selectbox(
                t("select_article"),
                range(len(article_titles)),
                format_func=lambda x: article_titles[x],
                key="article_selector"
            )
            
            if selected_idx is not None:
                article = st.session_state.articles[selected_idx]
                st.session_state.selected_article = selected_idx
                st.session_state.current_text = f"{article['Title']}\n\n{article['Description']}"
                
                # Article preview
                with st.expander(t("preview"), expanded=False):
                    st.markdown(f"**{article['Title'][:50]}...**")
                    st.caption(f"📍 {article['Source']}")
        else:
            st.info(t("click_refresh"))
        
        # ========== URL FETCH SECTION ==========
        st.markdown(f"<p style='text-align:center; color:#888;'>{t('or_divider')}</p>", unsafe_allow_html=True)
        
        st.markdown(f"**{t('fetch_url')}**")
        url_input = st.text_input(
            t("enter_url"),
            placeholder=t("url_placeholder"),
            key="url_input",
            label_visibility="collapsed"
        )
        
        if st.button(t("load_url"), use_container_width=True, disabled=not url_input):
            with st.spinner(t("loading_url")):
                try:
                    response = requests.post(
                        f"{API_URL}/fetch-url",
                        json={"url": url_input},
                        timeout=30
                    )
                    if response.status_code == 200:
                        data = response.json()
                        title = data.get("title", "Untitled")
                        content = data.get("content", "")
                        
                        if content:
                            st.session_state.current_text = f"{title}\n\n{content}"
                            st.session_state.fetched_url = url_input
                            st.success(f"✅ {t('url_loaded')}")
                        else:
                            st.warning(t("url_error") + " (empty content)")
                    else:
                        st.error(f"{t('url_error')}: {response.json().get('error', 'Unknown error')}")
                except requests.exceptions.ConnectionError:
                    st.error(t("backend_error"))
                except requests.exceptions.Timeout:
                    st.error(t("timeout"))
                except Exception as e:
                    st.error(f"{t('error')} {e}")
    
    # ========== TEXT INPUT MODE ==========
    elif st.session_state.input_mode == "text":
        st.markdown(f"### {t('text_input')}")
        
        text_input = st.text_area(
            t("enter_text"),
            height=150,
            placeholder=t("paste_here"),
            key="sidebar_text_input"
        )
        
        if text_input:
            st.session_state.current_text = text_input
            st.caption(f"📊 {len(text_input.split())} {t('words')}")
        
        # Quick examples
        with st.expander(t("examples")):
            examples = {
                "Tech": "The API rate limiting kicked in after we exceeded the threshold.",
                "Business": "Let's circle back on this and take it offline.",
                "Idiom": "Rome wasn't built in a day.",
            }
            for label, example in examples.items():
                if st.button(f"📝 {label}", key=f"ex_{label}", use_container_width=True):
                    st.session_state.current_text = example
                    st.rerun()
    
    # ========== SETTINGS MODE ==========
    else:
        st.markdown(f"### {t('manage_feeds')}")
        
        # Initialize custom feeds in session state
        if "custom_feeds" not in st.session_state:
            st.session_state.custom_feeds = []
        if "editing_feed" not in st.session_state:
            st.session_state.editing_feed = None
        
        # Load custom feeds
        try:
            response = requests.get(f"{API_URL}/custom-feeds", timeout=5)
            if response.status_code == 200:
                feeds = response.json().get("feeds")
                st.session_state.custom_feeds = feeds if feeds else []
        except:
            st.session_state.custom_feeds = []
        
        # Add new feed section
        with st.expander(t("add_feed"), expanded=False):
            new_title = st.text_input(t("feed_title"), key="new_feed_title", placeholder="MIT Technology Review")
            new_url = st.text_input(t("feed_url"), key="new_feed_url", placeholder="https://example.com/feed.xml")
            
            categories = [t("category_english"), t("category_tech"), t("category_medical"), t("category_news"), t("category_other")]
            new_category = st.selectbox(t("feed_category"), categories, key="new_feed_category")
            
            if st.button(t("save_feed"), key="add_new_feed", use_container_width=True):
                if new_title and new_url:
                    try:
                        response = requests.post(
                            f"{API_URL}/custom-feeds",
                            json={"title": new_title, "url": new_url, "category": new_category, "enabled": True},
                            timeout=5
                        )
                        if response.status_code == 200:
                            st.success(t("feed_added"))
                            # Clear RSS sources cache to reload
                            st.session_state.rss_sources = []
                            st.rerun()
                        else:
                            st.error(f"{t('error')} {response.status_code}")
                    except Exception as e:
                        st.error(f"{t('error')} {e}")
                else:
                    st.warning(f"⚠️ {t('feed_title')} and {t('feed_url')} are required")
        
        st.divider()
        
        # Default feeds (from config)
        st.markdown(f"**{t('default_feeds')}**")
        if st.session_state.rss_sources:
            custom_feed_titles = [cf.get("title") for cf in (st.session_state.custom_feeds or [])]
            for src in st.session_state.rss_sources:
                if src.get("title") not in custom_feed_titles:
                    category = src.get("category", "")
                    st.caption(f"📌 {src['title']} ({category})")
        else:
            st.caption(t("no_sources"))
        
        st.divider()
        
        # Custom feeds
        st.markdown(f"**{t('custom_feeds')}**")
        custom_feeds_list = st.session_state.custom_feeds or []
        if custom_feeds_list:
            for feed in custom_feeds_list:
                col_feed, col_action = st.columns([3, 1])
                with col_feed:
                    status = "✅" if feed.get("enabled", True) else "❌"
                    st.markdown(f"{status} **{feed['title']}**")
                    st.caption(f"{feed.get('category', '')} | {feed['url'][:40]}...")
                with col_action:
                    if st.button(t("delete_feed"), key=f"del_{feed['id']}", use_container_width=True):
                        try:
                            response = requests.delete(f"{API_URL}/custom-feeds/{feed['id']}", timeout=5)
                            if response.status_code == 200:
                                st.success(t("feed_deleted"))
                                st.session_state.rss_sources = []
                                st.rerun()
                        except Exception as e:
                            st.error(f"{t('error')} {e}")
                st.divider()
        else:
            st.info(t("no_custom_feeds"))
    
    # Footer
    st.caption("---")
    st.caption(t("footer2"))
    st.caption(t("footer3"))

# ==================== MAIN CONTENT ====================
st.title(t("main_title"))
st.caption(t("main_subtitle"))

# ========== ACTION BUTTONS (only in article/text mode) ==========
if st.session_state.input_mode != "settings":
    current_text = st.session_state.current_text
    has_text = bool(current_text.strip()) if current_text else False

    # Initialize TTS state
    if "is_reading" not in st.session_state:
        st.session_state.is_reading = False
    
    # Action buttons - give enough width to prevent text wrapping
    cols = st.columns([1, 1.2, 1.2, 1.2, 1.2, 1.2, 1.2, 1.2, 1, 1])
    
    with cols[1]:
        if st.button(t("explain"), use_container_width=True, help=t("help_explain"), disabled=not has_text):
            process_task("explain", current_text)
    with cols[2]:
        if st.button(t("summarize"), use_container_width=True, help=t("help_summarize"), disabled=not has_text):
            process_task("summarize", current_text)
    with cols[3]:
        if st.button(t("translate"), use_container_width=True, help=t("help_translate"), disabled=not has_text):
            process_task("translate", current_text)
    with cols[4]:
        # Toggle Read/Stop button
        if st.session_state.is_reading:
            # Show Stop button
            if st.button(t("stop_reading"), use_container_width=True, help=t("help_stop"), type="primary"):
                st.session_state.tts_action = "stop"
                st.session_state.is_reading = False
        else:
            # Show Read button
            if st.button(t("read_article"), use_container_width=True, help=t("help_read"), disabled=not has_text):
                st.session_state.tts_text = current_text
                st.session_state.tts_action = "speak"
                st.session_state.is_reading = True
    with cols[5]:
        if st.button(t("refine"), use_container_width=True, help=t("help_refine"), disabled=not has_text):
            process_task("refine", current_text)
    with cols[6]:
        if st.button(t("sentences"), use_container_width=True, help=t("help_sentences"), disabled=not has_text):
            process_task("extract_sentences", current_text)
    with cols[7]:
        if st.button(t("vocabulary"), use_container_width=True, help=t("help_vocabulary"), disabled=not has_text):
            process_task("extract_vocabulary", current_text)
    with cols[8]:
        if st.button(t("clear"), use_container_width=True, disabled=not st.session_state.messages):
            st.session_state.messages = []
            st.rerun()
    
    # Execute TTS action if set
    if "tts_action" in st.session_state:
        if st.session_state.tts_action == "speak" and "tts_text" in st.session_state:
            speak_text(st.session_state.tts_text)
            del st.session_state.tts_action
            del st.session_state.tts_text
        elif st.session_state.tts_action == "stop":
            stop_speech()
            del st.session_state.tts_action

st.divider()

# Show current content
if st.session_state.current_text:
    with st.expander(t("current_text"), expanded=True):
        # Allow HTML rendering for RSS content that contains HTML tags
        st.markdown(st.session_state.current_text, unsafe_allow_html=True)
        
        # Extract and show links
        links = extract_links(st.session_state.current_text)
        if links:
            st.divider()
            st.markdown("### 🔗 Links Found")
            for i, (text, url) in enumerate(links):
                col_link, col_btn = st.columns([5, 1])
                with col_link:
                    st.markdown(f"**{text}**")
                    st.caption(url)
                with col_btn:
                    if st.button("📥 Fetch", key=f"fetch_link_{i}", help=f"Fetch content from {url}"):
                        with st.spinner(t("loading_url")):
                            try:
                                response = requests.post(
                                    f"{API_URL}/fetch-url",
                                    json={"url": url},
                                    timeout=30
                                )
                                if response.status_code == 200:
                                    data = response.json()
                                    title = data.get("title", "Untitled")
                                    content = data.get("content", "")
                                    
                                    if content:
                                        st.session_state.current_text = f"{title}\n\n{content}"
                                        st.session_state.fetched_url = url
                                        st.success(f"✅ {t('url_loaded')}")
                                        st.rerun()
                                    else:
                                        st.warning(t("url_error") + " (empty content)")
                                else:
                                    st.error(f"{t('url_error')}: {response.json().get('error', 'Unknown error')}")
                            except Exception as e:
                                st.error(f"{t('error')} {e}")
else:
    if st.session_state.input_mode != "settings":
        st.info(t("use_sidebar"))

# Handle streaming response
if st.session_state.streaming_task:
    task_key = st.session_state.streaming_task
    text = st.session_state.streaming_text
    timestamp = st.session_state.streaming_timestamp
    
    task = TASKS[task_key]
    lang = st.session_state.get("language", "en")
    label = task["label_zh"] if lang == "zh" else task["label"]
    
    st.markdown(f"### ⚡ {t('processing')}...")
    
    # Show streaming response
    with st.chat_message("assistant", avatar="🤖"):
        response_placeholder = st.empty()
        full_response = ""
        
        try:
            for chunk in call_agent_stream(text, task_key):
                full_response += chunk
                response_placeholder.markdown(full_response + "▌")
            
            # Final display without cursor
            response_placeholder.markdown(full_response)
            
            # Save to messages
            st.session_state.messages.append({
                "role": "assistant", 
                "content": full_response,
                "timestamp": timestamp,
                "task": task_key
            })
        except Exception as e:
            st.error(f"Streaming error: {e}")
        finally:
            # Clear streaming state
            st.session_state.streaming_task = None
            st.session_state.streaming_text = None
            st.session_state.streaming_timestamp = None
            st.rerun()

# Results section
st.markdown(f"### {t('results')}")

if not st.session_state.messages:
    # Welcome message with translations
    lang = st.session_state.get("language", "en")
    st.markdown(f"""
    <div style="text-align: center; padding: 60px; color: #888; background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%); border-radius: 15px;">
        <h2>{t('welcome')}</h2>
        <p style="font-size: 18px;">{t('welcome_msg')}</p>
        <br>
        <p><strong>{t('quick_start')}</strong></p>
        <p>{t('step1')}</p>
        <p>{t('step2')}</p>
        <p>{t('step3')}</p>
        <br>
        <p><strong>{t('recommended')}</strong> {t('flow')}</p>
        <br>
        <p style="color: #666;">{t('tip')}</p>
    </div>
    """, unsafe_allow_html=True)
else:
    # Group messages into pairs (user + assistant) and reverse order (latest first)
    message_pairs = []
    for i in range(0, len(st.session_state.messages), 2):
        if i + 1 < len(st.session_state.messages):
            message_pairs.append((st.session_state.messages[i], st.session_state.messages[i + 1]))
        else:
            message_pairs.append((st.session_state.messages[i], None))
    
    # Display in reverse order (latest first)
    for idx, (user_msg, assistant_msg) in enumerate(reversed(message_pairs)):
        # Get timestamp
        timestamp = user_msg.get("timestamp", "")
        task_key = assistant_msg.get("task", "") if assistant_msg else ""
        
        # Display header with timestamp on the right
        col_title, col_time = st.columns([4, 1])
        with col_title:
            st.markdown(f"**{user_msg['content']}**")
        with col_time:
            st.markdown(f"<div style='text-align: right; color: #888; font-size: 0.85em;'>🕐 {timestamp}</div>", unsafe_allow_html=True)
        
        # Display assistant response with TTS buttons for sentences
        if assistant_msg:
            with st.chat_message("assistant", avatar="🤖"):
                render_content_with_tts(assistant_msg["content"], task_key, idx)
        
        st.divider()

# Footer
st.divider()
col_f1, col_f2, col_f3 = st.columns(3)
with col_f1:
    st.caption(t("footer1"))
with col_f2:
    st.caption(t("footer2"))
with col_f3:
    st.caption(t("footer3"))

# Execute TTS action if triggered from sentence buttons
if "tts_action" in st.session_state:
    if st.session_state.tts_action == "speak" and "tts_text" in st.session_state:
        speak_text(st.session_state.tts_text)
        del st.session_state.tts_action
        del st.session_state.tts_text
    elif st.session_state.tts_action == "stop":
        stop_speech()
        del st.session_state.tts_action
