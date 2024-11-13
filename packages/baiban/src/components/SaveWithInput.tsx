import React, { useState } from "react";

interface SaveButtonProps {
    onSave: (title: string) => void;
    initialTitle?: string; // 添加初始标题属性
}

export function SaveWithInput({ onSave, initialTitle = '未命名画布' }: SaveButtonProps) {
    const [isEditing, setIsEditing] = useState(false);
    const [titleInput, setTitleInput] = useState<string>(initialTitle); // 使用initialTitle初始化

    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setTitleInput(event.target.value);
    };

    const handleClickTitle = () => {
        setIsEditing(true);
    };

    const handleBlur = () => {
        if (titleInput.trim() === '') {
            alert('请输入标题！');
        } else {
            setIsEditing(false);
            onSave(titleInput);
        }
    };

    return (
        <div className="flex items-center bottom-12 right-4 space-x-2">
            {isEditing ? (
                // 当处于编辑状态时，显示输入框
                <input
                    type="text"
                    value={titleInput}
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
                    {titleInput}
                </span>
            )}
        </div>
    );
}