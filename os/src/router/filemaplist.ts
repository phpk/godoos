export const fileMapList = [
    {
        name: "document",
        icon: "word",
        editor: "/docx/index.html",
        //editor:"http://localhost:3001/",
        ext: ["docx", "doc",'docd'],
        hasPrview: true,
        eventType: "exportDocx"
    },
    {
        name: "excel",
        icon: "excel",
        editor: "/excel/index.html",
        ext: ["xlsx", "xls",'xlsd'],
        hasPrview: false,
        eventType: "exportExcel"
    },
    {
        name: "markdown",
        icon: "markdown",
        hasPrview: true,
        //priview:  "MarkDown",
        editor: "/markdown/index.html",
        eventType: "exportMd",
        ext: ['md']
    },
    {
        name: "mindmap",
        icon: "mindexe",
        editor: "/mind/index.html",
        hasPrview: false,
        //editor:"http://localhost:3003/",
        ext: ['mind'],
        eventType: "exportMind"
    },
    {
        name: "ppt",
        icon: "pptexe",
        editor: "/ppt/index.html",
        //editor:"http://localhost:5173/",
        ext: ['pptx', 'ppt', 'pptd'],
        hasPrview: false,
        eventType: "exportPPTX"
    },
    {
        name: 'fileEditor',
        icon: "editorbt",
        editor: "/text/index.html",
        hasPrview: false,
        ext: ['txt', 'html', 'json', 'xml', 'css', 'js', 'vue', 'go', 'php', 'java', 'py'],
        eventType: "exportText"
    },
    {
        name: "board",
        icon: "kanban",
        editor: "/kanban/index.html",
        hasPrview: false,
        ext: ['kb'],
        eventType: "exportKanban"
    },
    {
        name: 'whiteBoard',
        icon: "baiban",
        editor: "/baiban/index.html",
        hasPrview: false,
        ext: ['bb'],
        eventType: "exportBaiban"
    },
    {
        name: 'piceditor',
        icon: "picedit",
        editor: "/paint/index.html",
        //editor: "http://localhost:8080/",
        ext: ['pic', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff'],
        hasPrview: false,
        eventType: "exportPhoto"
    },
    {
        name: 'gantt',
        icon: "gant",
        editor: "/gantt/index.html",
        ext: ['gant'],
        hasPrview: false,
        eventType: "exportGant"
    },
    // {
    //     name: "PDF",
    //     icon: "pdf",
    //     priview:  "PdfViewer",
    //     width: 800,
    //     height: 600,
    //     ext: ['pdf']
    // },
    // {
    //     name: "music",
    //     icon: "music",
    //     priview: "MusicViewer",
    //     width: 300,
    //     height: 220,
    //     ext: ['mp3']
    // },
    // {
    //     name: "gallery",
    //     icon: "gallery",
    //     priview:  "ImageViewer",
    //     width: 800,
    //     height: 600,
    //     ext: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff'],
    // },
    // {
    //     name: "video",
    //     icon: "video",
    //     priview:  "VideoViewer",
    //     width: 800,
    //     height: 600,
    //     ext: ['mp4']
    // },
    // {
    //     name: 'browser',
    //     icon: "brower",
    //     priview:  "Browser",
    //     width: 800,
    //     height: 600,
    //     ext: ['link']
    // },
];
export const getExportType = (exportType:string) => {
    return fileMapList.find(item => item.eventType === exportType);
}
export const getFileType = (fileType:string) => {
    return fileMapList.find(item => item.ext.includes(fileType));
}
export const getFileName = (filename:string) => {
    return fileMapList.find(item => item.name === filename);
}