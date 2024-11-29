let currentZIndex = 1000; // 初始 z-index 值

// 动态生成模态对话框的 HTML 结构
function createModal() {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.id = `modal-${currentZIndex}`; // 添加唯一 ID

    const modalContent = document.createElement('div');
    modalContent.className = 'modal-content';

    const closeBtn = document.createElement('span');
    closeBtn.className = 'modal-close';
    closeBtn.innerHTML = '&times;';

    const title = document.createElement('h2');
    title.id = 'modal-title';

    const contentContainer = document.createElement('div');
    contentContainer.id = 'modal-content-container';

    const btnGroup = document.createElement('div');
    btnGroup.className = 'modal-button-group';
    const okButton = document.createElement('button');
    okButton.id = 'ok-button';
    okButton.className = 'modal-button modal-button-ok';
    okButton.textContent = '确定';

    const cancelButton = document.createElement('button');
    cancelButton.id = 'cancel-button';
    cancelButton.className = 'modal-button modal-button-cancel';
    cancelButton.textContent = '取消';

    modalContent.appendChild(closeBtn);
    modalContent.appendChild(title);
    modalContent.appendChild(contentContainer);
    btnGroup.appendChild(okButton);
    btnGroup.appendChild(cancelButton);
    modalContent.appendChild(btnGroup);
    modal.appendChild(modalContent);

    // 设置初始 z-index
    modal.style.zIndex = currentZIndex;

    return modal;
}

// 显示模态对话框
function showModal(options) {
    const modal = createModal();
    document.body.appendChild(modal);

    const title = modal.querySelector('#modal-title');
    const contentContainer = modal.querySelector('#modal-content-container');
    const okButton = modal.querySelector('#ok-button');
    const cancelButton = modal.querySelector('#cancel-button');

    title.textContent = options.titleText || '标题';

    // 清空旧的内容
    while (contentContainer.firstChild) {
        contentContainer.removeChild(contentContainer.firstChild);
    }

    // 插入内容或自定义 HTML
    if (options.contentText) {
        const content = document.createElement('p');
        content.textContent = options.contentText;
        contentContainer.appendChild(content);
    } else if (options.contentHtml) {
        contentContainer.innerHTML = options.contentHtml;
    }

    if (!options.showOk) {
        okButton.style.display = "none";
    } else {
        okButton.style.display = "block";
    }

    if (!options.showCancel) {
        cancelButton.style.display = "none";
    } else {
        cancelButton.style.display = "block";
    }

    // 关闭模态对话框
    function closeModal() {
        modal.remove();
    }

    // 当用户点击关闭按钮时
    const closeBtn = modal.querySelector('.modal-close');
    closeBtn.onclick = function() {
        closeModal();
    }

    // 当用户点击确定按钮时
    okButton.onclick = function() {
        if (typeof options.onOkClick === 'function') {
            const shouldClose = options.onOkClick();
            if (shouldClose !== false) {
                closeModal();
            }
        } else {
            closeModal();
        }
    }

    // 当用户点击取消按钮时
    cancelButton.onclick = function() {
        if (typeof options.onCancelClick === 'function') {
            options.onCancelClick();
        }
        closeModal();
    }

    // 当用户点击模态对话框外部时
    modal.onclick = function(event) {
        if (event.target === modal) {
            closeModal();
        }
    }

    // 更新当前最高的 z-index
    currentZIndex++;
}

// 示例: 调用 showModal 函数并传递确定和取消按钮的回调函数
// showModal({
//     titleText: '提示',
//     content: '你确定要继续吗?',
//     showOk: true,
//     showCancel: true,
//     onOkClick: function() {
//         alert('您点击了确定按钮');
//     },
//     onCancelClick: function() {
//         alert('您点击了取消按钮');
//     }
// });

// 创建第二个模态框
// setTimeout(function() {
//     showModal({
//         titleText: '回复',
//         contentHtml: '<textarea id="reply" style="width:100%;height:200px;"></textarea>',
//         showOk: true,
//         onOkClick: function() {
//             var reply = document.getElementById('reply').value;
//             if (reply == '') {
//                 showModal({
//                     titleText: '提示',
//                     content: '回复内容不能为空',
//                     showOk: true
//                 })
//                 return false; // 阻止关闭模态框
//             }
//         }
//     });
// }, 2000);