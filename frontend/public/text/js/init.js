window.Editor = {
    appearance: () => {
        var parent_window = (window.self == window.top),
            standalone = (window.matchMedia("(display-mode: standalone)").matches || navigator.standalone || !window.menubar.visible),
            apple_home_screen = (/(macOS|Mac|iPhone|iPad|iPod)/i.test(("userAgentData" in navigator) ? navigator.userAgentData.platform : navigator.platform) && "standalone" in navigator && navigator.standalone),
            hidden_chrome = (window.outerWidth == window.innerWidth && window.outerHeight == window.innerHeight),
            window_controls_overlay = ((("windowControlsOverlay" in navigator) ? navigator.windowControlsOverlay.visible : false) && standalone),
            refresh_window_controls_overlay = () => {
                var visibility = window_controls_overlay, styling = document.documentElement.classList.contains("window-controls-overlay");
                if (visibility != styling) (visibility) ? document.documentElement.classList.add("window-controls-overlay") : document.documentElement.classList.remove("window-controls-overlay");
            },
            fullscreen = (window.matchMedia("(display-mode: fullscreen)").matches || (!window.screenY && !window.screenTop && hidden_chrome) || (!window.screenY && !window.screenTop && standalone)),
            safe_area_insets = {
                left: getSafeAreaInset("left"),
                right: getSafeAreaInset("right"),
                top: getSafeAreaInset("top"),
                bottom: getSafeAreaInset("bottom"),
            },
            titlebar_area_insets = {
                top: getTitlebarAreaInset("top"),
                width_left: getTitlebarAreaInset("width-left"),
                width_right: getTitlebarAreaInset("width-right"),
                height: getTitlebarAreaInset("height"),
            },
            device_pixel_ratio = window.devicePixelRatio.toFixed(2),
            refresh_device_pixel_ratio = () => {
                var ratio = Editor.appearance().device_pixel_ratio, styling = getRootStyleProperty("--device-pixel-ratio");
                if (ratio != styling) document.documentElement.style.setProperty("--device-pixel-ratio", ratio);
            };
        function getSafeAreaInset(section) {
            return parseInt(getRootStyleProperty(`--safe-area-inset-${section}`), 10);
        }
        function getTitlebarAreaInset(section) {
            var value = getRootStyleProperty(`--titlebar-area-inset-${section}`);
            if (value.includes("calc")) value = Function(`"use strict"; return (${value.replace(/calc/g, "").replace(/100vw/g, `${window.innerWidth}px`).replace(/px/g, "")})`)();
            return parseInt(value, 10);
        }
        function getRootStyleProperty(template) {
            return window.getComputedStyle(document.documentElement).getPropertyValue(template);
        }
        return { parent_window, standalone, apple_home_screen, hidden_chrome, window_controls_overlay, refresh_window_controls_overlay, fullscreen, safe_area_insets, titlebar_area_insets, device_pixel_ratio, refresh_device_pixel_ratio };
    },
    environment: () => ({
        file_protocol: (window.location.protocol == "file:"),
        touch_device: ("ontouchstart" in window),
        apple_device: (/(macOS|Mac|iPhone|iPad|iPod)/i.test(("userAgentData" in navigator) ? navigator.userAgentData.platform : navigator.platform)),
        macOS_device: (/(macOS|Mac)/i.test(("userAgentData" in navigator) ? navigator.userAgentData.platform : navigator.platform) && navigator.standalone == undefined),
        mozilla_browser: (CSS.supports("-moz-appearance: none"))
    }),
    support: () => ({
        local_storage: (window.location.protocol != "blob:") ? window.localStorage : null,
        file_system: ("showOpenFilePicker" in window),
        file_handling: ("launchQueue" in window && "LaunchParams" in window),
        window_controls_overlay: ("windowControlsOverlay" in navigator),
        editing_commands: (!Editor.environment().mozilla_browser),
        web_sharing: ("share" in navigator)
    }),
    query: (identifier = Editor.active_editor) => {
        var tab = workspace_tabs.querySelector(`.tab[data-editor-identifier="${identifier}"]`),
            container = workspace_editors.querySelector(`.editor[data-editor-identifier="${identifier}"]`),
            textarea = (container) ? container.editor : null,
            getName = section => {
                if ((document.querySelectorAll(`[data-editor-identifier="${identifier}"]:not([data-editor-change])`).length == 0) && (identifier != Editor.active_editor)) return null;
                var name = workspace_tabs.querySelector(`.tab[data-editor-identifier="${identifier}"] [data-editor-name]`).innerText;
                if (!section || (!name.includes(".") && section == "base")) return name;
                if (section == "base") {
                    name = name.split(".");
                    name.pop();
                    return name.join(".");
                }
                if (section == "extension") {
                    if (!name.includes(".")) return "";
                    return name.split(".").pop();
                }
            }
        return { tab, container, textarea, getName };
    },
    view: () => document.body.getAttribute("data-view"),
    view_change: () => (document.body.hasAttribute("data-view-change")),
    orientation: () => document.body.getAttribute("data-orientation"),
    orientation_change: () => (document.body.hasAttribute("data-orientation-change")),
    scaling_change: () => (document.body.hasAttribute("data-scaling-change")),
    unsaved_work: () => (!Editor.appearance().parent_window || (workspace_tabs.querySelectorAll(".tab:not([data-editor-change])[data-editor-unsaved]").length == 0)),
    preapproved_extensions: ["txt", "html", "css", "js", "php", "json", "webmanifest", "bbmodel", "xml", "yaml", "yml", "dist", "config", "ini", "md", "markdown", "mcmeta", "lang", "properties", "uidx", "material", "h", "fragment", "vertex", "fxh", "hlsl", "ihlsl", "svg"],
    active_editor: null,
    preview_editor: "active-editor",
    file_handles: {},
    child_windows: [],
    settings: {
        entries: JSON.parse(window.localStorage.getItem("settings")) || {},
        set: (key, value) => {
            if (!Editor.support().local_storage) return;
            Editor.settings.entries[key] = value;
            window.localStorage.setItem("settings", JSON.stringify(Editor.settings.entries, null, "  "));
            return value;
        },
        remove: key => {
            if (!Editor.support().local_storage) return;
            delete Editor.settings.entries[key];
            window.localStorage.setItem("settings", JSON.stringify(Editor.settings.entries, null, "  "));
            return true;
        },
        has: key => (key in Editor.settings.entries),
        get: key => {
            if (!Editor.support().local_storage) return;
            if (!Editor.settings.has(key)) return;
            return Editor.settings.entries[key];
        },
        reset: ({ confirm: showPrompt = false } = {}) => {
            if (!Editor.support().local_storage) return;
            if (showPrompt) {
                if (!confirm("确定要重置所有设置?")) return false;
            }
            default_orientation_setting.select("horizontal");
            setSyntaxHighlighting(false);
            syntax_highlighting_setting.checked = false;
            automatic_refresh_setting.checked = true;
            preview_base_input.reset();
            Editor.settings.entries = {};
            window.localStorage.removeItem("settings");
            if (showPrompt) reset_settings_card.open();
            return true;
        }
    },
    active_dialog: null,
    dialog_previous: null,
    active_widget: null,
    picker_color: null,
    install_prompt: null
};
if (Editor.appearance().parent_window) document.documentElement.classList.add("startup-fade");
if (Editor.appearance().apple_home_screen) document.documentElement.classList.add("apple-home-screen");
if (Editor.environment().touch_device) document.documentElement.classList.add("touch-device");
if (Editor.environment().apple_device) document.documentElement.classList.add("apple-device");
if (Editor.environment().macOS_device) document.documentElement.classList.add("macOS-device");
if (Editor.support().web_sharing) document.documentElement.classList.add("web-sharing");
Editor.appearance().refresh_window_controls_overlay();
Editor.appearance().refresh_device_pixel_ratio();
