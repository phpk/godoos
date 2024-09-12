export const appList = [
    {
        name: 'setting',
        appIcon: "setting",
        content: "Setting",
        multiple: false,
        width: 800,
        height: 600,
        frame: false,
        backgroundColor: '#00000000',
        center: true,
        resizable: true,
        isDeskTop: false,
        isMenuList: true
    },
    {
        name: 'computer',
        appIcon: "diannao",
        content: "Computer",
        width: 800,
        height: 600,
        frame: true,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMenuList: true,
        dir: '/',
    },

    {
        name: 'appstore',
        multiple: false,
        appIcon: "store",
        content: "Store",
        frame: true,
        width: 900,
        height: 600,
        center: true,
        resizable: true,
        backgroundColor: '#ffffff00',
        isDeskTop: true,
        isMenuList: true,
    },
    {
        name: 'system.version',
        appIcon: "info",
        content: "Version",
        width: 300,
        height: 200,
        frame: true,
        center: true,
        resizable: false,
        isDeskTop: false,
        isMenuList: true,
    },
    {
        name: 'process.title',
        appIcon: "progress",
        content: "ProcessManager",
        width: 800,
        height: 600,
        frame: true,
        center: true,
        resizable: false,
        isDeskTop: false,
        isMenuList: true,
    },
    
    {
        name: "create.shortcut",
        appIcon: "link",
        content: "CreateUrl",
        width: 400,
        height: 400,
        resizable: true,
        center: true,
        isDeskTop: false,
        isContext: true
    },
    {
        name: "document",
        appIcon: "word",
        url: "/docx/index.html",
        width: 800,
        frame: true,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMagnet: false,
        ext: ["docx", "doc"],
        eventType: "exportDocx"
    },
    {
        name: "excel",
        appIcon: "excel",
        url: "/excel/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMagnet: false,
        ext: ["xlsx", "xls"],
        eventType: "exportExcel"
    },
    {
        name: "markdown",
        appIcon: "markdown",
        content: "MarkDown",
        width: 800,
        height: 600,
        frame: true,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMagnet: false,
        ext: ['md']
    },
    {
        name: "mindmap",
        appIcon: "mindexe",
        url: "/mind/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMagnet: false,
        ext: ['mind'],
        eventType: "exportMind"
    },
    {
        name: "ppt",
        appIcon: "pptexe",
        url: "/ppt/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMagnet: false,
        ext: ['pptx', 'ppt'],
        eventType: "exportPPTX"
    },
    {
        name: 'fileEditor',
        appIcon: "editorbt",
        url: "/text/index.html",
        width: 800,
        height: 600,
        frame: true,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMagnet: false,
        ext: ['txt', 'html', 'json', 'xml', 'css', 'js', 'vue', 'go', 'php', 'java', 'py'],
        eventType: "exportText"
    },
    {
        name: "board",
        appIcon: "kanban",
        url: "/kanban/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        ext: ['kb'],
        eventType: "exportKanban"
    },
    {
        name: 'whiteBoard',
        appIcon: "baiban",
        url: "/baiban/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        ext: ['bb'],
        eventType: "exportBaiban"
    },
    {
        name: 'localchat',
        appIcon: "chat",
        //content: "LocalChat",
        content:"Chat",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
    },

    {
        name: 'piceditor',
        appIcon: "picedit",
        url: "/picedit/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        ext: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff'],
        eventType: "exportPhoto"
    },
    {
        name: 'gantt',
        appIcon: "gant",
        url: "/gantt/index.html",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        ext: ['gant'],
        eventType: "exportGant"
    },
    {
        name: 'calculator',
        appIcon: "calculator",
        url: "/calculator/index.html",
        frame: true,
        width: 366,
        height: 550,
        center: true,
        resizable: false,
        isDeskTop: false,
        isMenuList: true,
    },
    {
        name: 'calendar',
        appIcon: "calendar",
        content: "Calendar",
        frame: true,
        width: 562,
        height: 495,
        center: true,
        resizable: true,
        isDeskTop: false,
        isMenuList: true,
    },
    {
        name: "PDF",
        appIcon: "pdf",
        content: "PdfViewer",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: false,
        isMaximize: true,
        ext: ['pdf']
    },
    {
        name: "musicStore",
        appIcon: "music",
        content: "MusicStore",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: false,
        isMenuList: true,
    },
    {
        name: "music",
        appIcon: "music",
        view: "MusicViewer",
        frame: true,
        width: 300,
        height: 220,
        center: true,
        resizable: true,
        isDeskTop: false,
        isMenuList: false,
        ext: ['mp3']
    },
    {
        name: "gallery",
        appIcon: "gallery",
        content: "PictureStore",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: false,
        isMenuList: true,
    },
    {
        name: "gallery",
        appIcon: "gallery",
        content: "ImageViewer",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: false,
        isMenuList: false,
        ext: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff'],
    },
    {
        name: "video",
        appIcon: "video",
        content: "VideoViewer",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        isDeskTop: false,
        isMenuList: false,
        ext: ['mp4']
    },
    {
        name: 'browser',
        appIcon: "brower",
        content: "Browser",
        frame: true,
        width: 800,
        height: 600,
        center: true,
        resizable: true,
        isDeskTop: true,
        isMenuList: true,
    },
];