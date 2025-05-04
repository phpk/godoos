// src/router/index.ts
import { useWindowStore } from '@/stores/window'
import { isMobileDevice } from '@/utils/device'
import { markRaw } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../App.vue'),
  },
  {
    path: '/setting',
    name: 'setting',
    component: () => import('../views/setting.vue'),
    meta: {
      title: 'setting',
      icon: 'setting',
      windowConfig: {
        width: 800,
        height: 600,
      },
    },
    props: {
      isOne: true,
    },
  },
  {
    path: '/computer',
    name: 'FileExplorer',
    component: () => import('../views/computer.vue'),
    meta: {
      title: 'computer',
      icon: 'diannao',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/recycle',
    name: 'recycle',
    component: () => import('../views/computer.vue'),
    meta: {
      title: 'recycle',
      icon: 'trush',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      path: '/B',
    },
  },
  {
    path: '/viewer',
    name: 'viewer',
    component: () => import('../views/viewer.vue'),
    meta: {
      title: 'viewer',
      icon: 'viewer',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      isFile: true,
      action: 'preview',
    },
  },
  {
    path: '/editor',
    name: 'editor',
    component: () => import('../views/editor.vue'),
    meta: {
      title: 'editor',
      icon: 'editor',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      isFile: true,
      action: 'editor',
    },
  },
  {
    path: '/document',
    name: 'document',
    component: {
      template: '/docx/index.html',
    },
    meta: {
      title: 'document',
      icon: 'word',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/fileEditor',
    name: 'fileEditor',
    component: {
      template: '/text/index.html',
    },
    meta: {
      title: 'fileEditor',
      icon: 'editorbt',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/excel',
    name: 'excel',
    component: {
      template: '/excel/index.html',
    },
    meta: {
      title: 'excel',
      icon: 'excel',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/markdown',
    name: 'markdown',
    component: {
      template: '/markdown/index.html',
    },
    meta: {
      title: 'markdown',
      icon: 'markdown',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/mindmap',
    name: 'mindmap',
    component: {
      template: '/mind/index.html',
    },
    meta: {
      title: 'mindmap',
      icon: 'mindexe',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/ppt',
    name: 'ppt',
    component: {
      template: '/ppt/index.html',
    },
    meta: {
      title: 'ppt',
      icon: 'pptexe',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/board',
    name: 'board',
    component: {
      template: '/kanban/index.html',
    },
    meta: {
      title: 'board',
      icon: 'kanban',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/whiteBoard',
    name: 'whiteBoard',
    component: {
      template: '/baiban/index.html',
    },
    meta: {
      title: 'whiteBoard',
      icon: 'baiban',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/gantt',
    name: 'gantt',
    component: {
      template: '/gantt/index.html',
    },
    meta: {
      title: 'gantt',
      icon: 'gant',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/calculator',
    name: 'calculator',
    component: {
      template: '/calculator/index.html',
    },
    meta: {
      title: 'calculator',
      icon: 'calculator',
      windowConfig: {
        width: 420,
        height: 580,
      },
    },
  },
  {
    path: '/calendar',
    name: 'calendar',
    component: () => import('../components/taskbar/Calendar.vue'),
    meta: {
      title: 'calendar',
      icon: 'calendar',
      windowConfig: {
        width: 562,
        height: 518,
      },
    },
  },
  {
    path: '/piceditor',
    name: 'piceditor',
    component: {
      template: '/paint/index.html',
    },
    meta: {
      title: 'piceditor',
      icon: 'picedit',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
  {
    path: '/workchat',
    name: 'workbench',
    component: () => import('../views/workbench.vue'),
    meta: {
      title: 'workchat',
      icon: 'workchat',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      isOne: true,
    },
  },
  {
    path: '/version',
    name: 'version',
    component: () => import('../views/version.vue'),
    meta: {
      title: 'version',
      icon: 'info',
      windowConfig: {
        width: 300,
        height: 200,
      },
    },
    props: {
      isOne: true,
    },
  },
  {
    path: '/assistant',
    name: 'assistant',
    component: () => import('../views/assistant.vue'),
    meta: {
      title: 'assistant',
      icon: 'aiassistant',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      isOne: true,
    },
  },
  {
    path: '/aiModule',
    name: 'aiModule',
    component: () => import('../views/aiModule.vue'),
    meta: {
      title: 'aiModule',
      icon: 'aidown',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      isOne: true,
    },
  },
  {
    path: '/aiSetting',
    name: 'aiSetting',
    component: () => import('../views/aiSetting.vue'),
    meta: {
      title: 'aiSetting',
      icon: 'aisetting',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
    props: {
      isOne: true,
    },
  },
  {
    path: '/knowledgechat',
    name: 'knowledgechat',
    component: () => import('../views/knowledgeChat.vue'),
    meta: {
      title: '知识库对话',
      icon: 'aichat',
      windowConfig: {
        width: 900,
        height: 600,
      },
    },
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

// 路由守卫，用于在窗口中打开路由页面
router.beforeEach(async (to, _from, next) => {
  // console.log('Route guard triggered'); // 添加调试信息
  console.log('to:', to)
  // console.log('from:', from);
  // console.log('next:', next);

  if (to.name && to.meta && to.meta.icon) {
    const windowStore = useWindowStore()
    const matched: any = to.matched[0]
    const defaultComponent = matched.components.default
    if (matched && defaultComponent) {
      let component: any

      if (
        defaultComponent.template &&
        typeof defaultComponent.template === 'string'
      ) {
        component = defaultComponent.template
      } else {
        if (to.query.editor && to.query.hasPrview == 'false') {
          to.query.action = 'edit'
          component = to.query.editor
        } else {
          const componentPromise = defaultComponent as () => Promise<any>
          component = await componentPromise()
          component = markRaw(component.default || component)
        }
      }
      //console.log('component:', component)
      const windowConfig = to.meta.windowConfig as
        | { width: number; height: number }
        | undefined

      if (windowConfig) {
        //console.log(to)
        windowStore.create({
          title: to.meta.title as string,
          icon: to.meta.icon as string,
          component: component, // 使用解析后的组件
          props: {
            ...to.params,
            ...matched.props.default,
            ...to.query,
          },
          size: windowConfig,
          isMaximized: isMobileDevice() ? true : false,
        })
      }
      return next('/')
    }
  } else {
    next()
  }
})

//console.log('Router guard registered'); // 添加调试信息

export default router
