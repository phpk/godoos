import React from 'react';

// 假设您的SVG图片位于public/images/save.svg
export const SaveButton = ({ onSave }) => {
  const handleClick = () => {
    if (onSave) {
      onSave();
    }
  };

  return (
    <button
      type="button"
      className="bg-transparent border border-gray-300 mr-3 rounded-lg p-1 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-600"
      onClick={handleClick}
    >
      {/* 使用img标签直接引用SVG图片 */}
      <img
        src="./save.svg" // 注意这里的路径是相对于public目录的
        alt="Save"
        className="w-6 h-6 inline-block align-middle"
      />
    </button>
  );
};
