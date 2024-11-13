import { Tldraw, createTLStore, defaultShapeUtils, throttle } from 'tldraw'
import { getAssetUrls } from '@tldraw/assets/selfHosted'
//import { SaveWithInput } from "./components/SaveWithInput";
import { SaveButton } from "./components/SaveButton";
import { useLayoutEffect, useState } from "react";
function isBase64(str) {
  if (str === '' || str.trim() === '') {
    return false
  }
  try {
    return btoa(atob(str)) == str
  } catch (err) {
    return false
  }
}
function decodeBase64(base64String) {
  // 将Base64字符串分成每64个字符一组
  const padding =
    base64String.length % 4 === 0 ? 0 : 4 - (base64String.length % 4)
  base64String += '='.repeat(padding)

  // 使用atob()函数解码Base64字符串
  const binaryString = atob(base64String)

  // 将二进制字符串转换为TypedArray
  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }

  // 将TypedArray转换为字符串
  return new TextDecoder('utf-8').decode(bytes)
}
const PERSISTENCE_KEY = 'PERSISTENCE_KEY'
export default function App() {
  const assetUrls = getAssetUrls({
    baseUrl: './static/',
  })
  localStorage.removeItem(PERSISTENCE_KEY)
  const [store] = useState(() => createTLStore({ shapeUtils: defaultShapeUtils }))
  //[2]
  const [loadingState, setLoadingState] = useState({
    status: 'loading',
  })
  const [isEditing, setIsEditing] = useState(false);
  const [canvasTitle, setCanvasTitle] = useState('未命名白板');
  const handleTitleChange = (event) => {
    setCanvasTitle(event.target.value);
  };

  const handleClickTitle = () => {
    setIsEditing(true);
  };

  const handleBlur = () => {
    if (canvasTitle.trim() === '') {
      alert('请输入标题！');
    } else {
      setIsEditing(false);
      //setCanvasTitle(canvasTitle.value);
    }
  };

  useLayoutEffect(() => {
    setLoadingState({ status: 'loading' })
    // Get persisted data from local storage
    const persistedSnapshot = localStorage.getItem(PERSISTENCE_KEY)

    if (persistedSnapshot) {
      try {
        const snapshot = JSON.parse(persistedSnapshot)
        store.loadSnapshot(snapshot)
        setLoadingState({ status: 'ready' })
      } catch (error) {
        setLoadingState({ status: 'error', error: error.message }) // Something went wrong
      }
    } else {
      setLoadingState({ status: 'ready' }) // Nothing persisted, continue with the empty store
    }
    const eventHandler = (e) => {
      try {
        //console.log(e)
        const eventData = e.data
        //console.log(eventData)
        if (eventData.type === 'start') {
          if (eventData.title) {
            //console.log(eventData.title)
            const baseTitle = eventData.title.substring(0, eventData.title.lastIndexOf('.'))
            //console.log(baseTitle)
            if (baseTitle) {
              setCanvasTitle(baseTitle)
            }

          }
        }
        if (eventData.type === 'init') {
          const data = eventData.data
          //console.log(data)
          if (!data || !data.title) {
            return;
          }
          const title = data.title.substring(0, data.title.lastIndexOf('.'))
          setCanvasTitle(title)

          if (data.content) {
            //instanceEditor.store.clear()
            if (typeof data.content === 'string' && isBase64(data.content)) {
              data.content = decodeBase64(data.content)
              data.content = JSON.parse(data.content)
            }
            //console.log(data)
            store.loadSnapshot(data.content)
            localStorage.setItem(PERSISTENCE_KEY, JSON.stringify(data.content))
          }

        }
      } catch (error) {
        console.log(error)
      }

    }
    // window.addEventListener('load', () => {
    //   window.parent.postMessage({ type: 'initSuccess' }, '*');
    //   window.addEventListener('message', eventHandler);
    // });
    //window.onload = function () {
      //console.log('addEventListener init')
      window.parent.postMessage({ type: 'initSuccess' }, '*');
      //console.log('addEventListener init')
      window.addEventListener('message', eventHandler);
    //};
    // Each time the store changes, run the (debounced) persist function
    const cleanupFn = store.listen(
      throttle(() => {
        //const snapshot = store.getSnapshot()
        const snapshot = store.getSnapshot()
        localStorage.setItem(PERSISTENCE_KEY, JSON.stringify(snapshot))
      }, 200)
    )

    return () => {
      cleanupFn();
      window.removeEventListener('unload', eventHandler);
    }
  }, [store]);



  const handleSave = () => {

    // 保存逻辑
    const content = localStorage.getItem(PERSISTENCE_KEY) || '{}'
    const saveData = {
      title: canvasTitle,
      content: JSON.parse(content)
    }
    console.log(saveData)
    const save = {
      data: JSON.stringify(saveData),
      type: 'exportBaiban'
    }

    window.parent.postMessage(save, '*')
  };
  if (loadingState.status === 'loading') {
    return (
      <div className="tldraw__editor">
        <h2>Loading...</h2>
      </div>
    )
  }

  if (loadingState.status === 'error') {
    return (
      <div className="tldraw__editor">
        <h2>Error!</h2>
        <p>{loadingState.error}</p>
      </div>
    )
  }
  return (
    <div className="tldraw__editor">
      <div style={{
        position: 'absolute',
        zIndex: 2000,
        top: 3,
        left: 'calc(50% - 90px)',
        display: 'flex', // 添加这一行启用Flex布局
        alignItems: 'center', // 可选，使按钮和输入框垂直居中对齐
      }}>
        <SaveButton onSave={() => handleSave()} />
        <div className="flex items-center bottom-12 right-4 space-x-2">
          {isEditing ? (
            // 当处于编辑状态时，显示输入框
            <input
              type="text"
              value={canvasTitle}
              onChange={handleTitleChange}
              onBlur={handleBlur}
              className="
                        bg-gray-100 
                        border 
                        border-gray-200 
                        rounded-lg 
                        py-1.5 
                        px-3 
                        w-64 
                        focus:outline-none 
                        focus:border-slate-500
                    "
              style={{ height: '38px' }}
            />
          ) : (
            // 默认显示未命名画布文字，点击时切换到编辑状态
            <span
              onClick={handleClickTitle}
              className="cursor-pointer select-none text-gray-400 hover:text-gray-600 transition duration-200 ease-in-out mt-2"
            >
              {canvasTitle}
            </span>
          )}
        </div>
      </div>
      <Tldraw
        assetUrls={assetUrls}
        persistenceKey="my-persistence-key"
        store={store}
      />
    </div>
  )
}