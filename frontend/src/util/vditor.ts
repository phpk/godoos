const defaultHeight = 350
export function getEditorTitleOption() {
  //const editorTheme: string = theme.global.current.value.dark ? "dark" : "classic";
  //const editorTheme: string = isdark ? "dark" : "classic";
  return {
    height: defaultHeight,
    cdn: "/vditor",
    theme: "classic",
    toolbarConfig: {
      pin: true,
    },
    counter: {
      enable: true
    },
    resize: {
      enable: true
    },
    toolbar: [
      "headings",
      
      "list",
      "ordered-list",
      
      {
        name: "more",
        toolbar: [
          "edit-mode",
          "check",
          "quote",
          "bold",
          "line",
          "code",
          "table",
          "italic",
          "strike",
          "link",
          "outdent",
          "indent",
          "inline-code",
          "insert-before",
          "insert-after",
          "emoji",
          "undo",
          "redo",
          "both",
          // "upload",
          // "record",
          "code-theme",
          "content-theme",
          "export",
          "outline",
          "preview",
          "fullscreen",
        ],
      },
    ],
    cache: {
      enable: false,
    },
    after: () => {
      //baseContent.value.setValue("");
    },
  };

}
export function getEditorOption() {
  //const editorTheme: string = theme.global.current.value.dark ? "dark" : "classic";
  //const editorTheme: string = isdark ? "dark" : "classic";
  return {
    height: defaultHeight,
    cdn: "/vditor",
    theme: "classic",
    toolbarConfig: {
      pin: true,
    },
    counter: {
      enable: true
    },
    resize: {
      enable: true
    },
    toolbar: [
      "headings",
      "bold",
      "italic",
      "strike",
      "link",
      "list",
      "ordered-list",
      "check",
      
      "table",
      "edit-mode",
      {
        name: "more",
        toolbar: [
          "quote",
          "line",
          "code",
          "outdent",
          "indent",
          "inline-code",
          "insert-before",
          "insert-after",
          "emoji",
          "undo",
          "redo",
          "both",
          // "upload",
          // "record",
          "code-theme",
          "content-theme",
          "export",
          "outline",
          "preview",
          "fullscreen",
        ],
      },
    ],
    cache: {
      enable: false,
    },
    after: () => {
      //baseContent.value.setValue("");
    },
  };

}
export function getMdOption() {
  //const editorTheme: string = theme.global.current.value.dark ? "dark" : "classic";
  //const editorTheme: string = isdark ? "dark" : "classic";
  return {
    height: 500,
    cdn: "/vditor",
    theme: "classic",
    toolbarConfig: {
      pin: true,
    },
    counter: {
      enable: true
    },
    resize: {
      enable: true
    },
    toolbar: [
      "edit-mode",
      "headings",
      "bold",
      "italic",
      "strike",
      "link",
      "list",
      "ordered-list",
      "check",
      "table",
      "edit-mode",
      "quote",
      "line",
      "code",
      "outdent",
      "indent",
      "inline-code",
      "insert-before",
      "insert-after",
      
      "undo",
      "redo",
      "both",
      "outline",
      "preview",
      
      {
        name: "more",
        toolbar: [
          "fullscreen",
          "code-theme",
          "emoji",
          "content-theme",
          "export",
        ],
      }
    ],
    cache: {
      enable: false,
    },
    after: () => {
      //baseContent.value.setValue("");
    },
  };

}


