// 导入必要的React库和类型
import React, { createContext, useState } from "react";

// 定义上下文的类型别名，清晰指定上下文中包含的数据结构
interface TitleContextType {
    title: string; // 当前的标题
    setTitle: (newTitle: string) => void; // 设置新标题的函数
}

// 创建一个上下文实例，初始化默认值
export const TitleContext = createContext<TitleContextType>({
    title: "未命名看板", // 默认标题
    setTitle: () => {}, // 默认的空函数，用于设置标题
});

// 定义一个组件，它作为上下文的提供者
export const TitleContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    // 使用useState Hook来管理标题状态
    // 使用useState Hook来管理标题状态
    const [title, setTitle] = useState<string>("未命名看板"); // 注意这里改名为setTitleRaw

    // 添加日志记录的setTitle函数
    // const setTitle = (newTitle: string) => {
    //     console.log("setTitle function executed with:", newTitle);
    //     setTitleRaw(newTitle);
    // };
    // 准备提供给上下文的值对象
    const contextValue = {
        title, // 当前的标题状态
        setTitle, // 设置标题的函数
    };

    // 使用TitleContext.Provider包裹子组件，使得其后代组件可以访问到title和setTitle
    return (
        <TitleContext.Provider value={contextValue}>
            {children} {/* 渲染传递进来的子组件 */}
        </TitleContext.Provider>
    );
};