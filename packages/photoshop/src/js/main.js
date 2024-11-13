/**
 * miniPaint - https://github.com/viliusle/miniPaint
 * author: Vilius L.
 */

//css
import './../css/reset.css';
import './../css/utility.css';
import './../css/component.css';
import './../css/layout.css';
import './../css/menu.css';
import './../css/print.css';
import './../../node_modules/alertifyjs/build/css/alertify.min.css';
//js
import app from './app.js';
import config from './config.js';
import './core/components/index.js';
import Base_gui_class from './core/base-gui.js';
import Base_layers_class from './core/base-layers.js';
import Base_tools_class from './core/base-tools.js';
import Base_state_class from './core/base-state.js';
import Base_search_class from './core/base-search.js';
import File_open_class from './modules/file/open.js';
import File_save_class from './modules/file/save.js';
import * as Actions from './actions/index.js';
import { md5 } from 'js-md5'

window.addEventListener('load', function (e) {
	// Initiate app
	var Layers = new Base_layers_class();
	var Base_tools = new Base_tools_class(true);
	var GUI = new Base_gui_class();
	var Base_state = new Base_state_class();
	var File_open = new File_open_class();
	var File_save = new File_save_class();
	var Base_search = new Base_search_class();

	// Register singletons in app module
	app.Actions = Actions;
	app.Config = config;
	app.FileOpen = File_open;
	app.FileSave = File_save;
	app.GUI = GUI;
	app.Layers = Layers;
	app.State = Base_state;
	app.Tools = Base_tools;

	// Register as global for quick or external access
	window.Layers = Layers;
	window.AppConfig = config;
	window.State = Base_state;
	window.FileOpen = File_open;
	window.FileSave = File_save;

	// Render all
	GUI.init();
	Layers.init();
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
	//godo
	const eventHandler = (e) => {
		const eventData = e.data
		console.log(eventData)
		if (eventData.type === 'init') {
			const data = eventData.data
			console.log(data)
			if (!data || !data.title) {
				return;
			}
			const ext = data.title.split(".").pop() || 'png'
			if(data.content) {
				if (typeof data.content === 'string' && isBase64(data.content)) {
					// data.content = decodeBase64(data.content)
					// data.content = JSON.parse(data.content)
					const fileName = 'photoshop_' + md5(data.content)
					const jsonData = localStorage.getItem(fileName)
					if (jsonData) {
						File_open.load_json(jsonData, false)
					} else {
						File_open.file_open_data_url_handler(`data:image/${ext};base64,` + data.content, false)
					}
				}
				if(data.content instanceof ArrayBuffer) {
					const blob = new Blob([data.content]);
					const reader = new FileReader();
					reader.onload = function (e) {
						//console.log(e.target.result)
						const readContent = e.target.result.split(',')[1];
						//console.log(readContent);
						const fileName = 'photoshop_' + md5(readContent)
						const jsonData = localStorage.getItem(fileName)
						//console.log(jsonData)
						if (jsonData) {
							File_open.load_json(jsonData, false)
						}else{
							File_open.file_open_data_url_handler(`data:image/${ext};base64,` + readContent, false)
						}

					};
					reader.readAsDataURL(blob);
				}
				// var img = new Image();
				// img.crossOrigin = "Anonymous";
				// img.onload = function () {
				// 	var new_layer = {
				// 		name: data.title,
				// 		type: 'image',
				// 		link: img,
				// 		width: img.width,
				// 		height: img.height,
				// 		width_original: img.width,
				// 		height_original: img.height,
				// 	};
				// 	Layers.insert(new_layer);
				// };
				// img.onerror = function (ex) {
				// 	alertify.error('Sorry, image could not be loaded. Try copy image and paste it.');
				// };
				// img.src = data.content;
			}
			//Layers.insert(new_layer);
		

		}
	}
	window.parent.postMessage({ type: 'initSuccess' }, '*')
	window.addEventListener('message', eventHandler);
	window.addEventListener('unload', () => {
		window.removeEventListener('message', eventHandler)
	})
	
}, false);
