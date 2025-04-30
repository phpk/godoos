<template>
    <div class="slides-content">
        <div class="slides-item slides-fade slides-bg1"></div>
        <div class="slides-item slides-fade slides-bg2"></div>
        <div class="slides-item slides-fade slides-bg3"></div>
    </div>
</template>
<style lang="scss" scoped>
.slides-content {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
}

.slides-item {
    position: absolute;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    opacity: 0;
    animation: fadeEffect 15s infinite;
    transition: opacity 0.3s ease;
}

.slides-bg1 {
    background-image: url('/public/images/bg/bg1.jpg');
    background-size: cover;
    background-position: center;
}

.slides-bg2 {
    background-image: url('/public/images/bg/bg2.jpg');
    background-size: cover;
    background-position: center;
}

.slides-bg3 {
    background-image: url('/public/images/bg/bg3.jpg');
    background-size: cover;
    background-position: center;
}

@keyframes fadeEffect {
    0% {
        opacity: 0;
    }

    5% {
        opacity: 1;
    }

    33.333% {
        opacity: 1;
    }

    38.333% {
        opacity: 0;
    }

    100% {
        opacity: 0;
    }
}

.slides-fade {
    animation-fill-mode: forwards;
    /* 当动画完成后，‌保持最后一个属性值 */
}

/* 为每张图片设置不同的动画延迟时间 */
.slides-item:nth-child(1) {
    animation-delay: 0s;
}

.slides-item:nth-child(2) {
    animation-delay: 5s;
}

.slides-item:nth-child(3) {
    animation-delay: 10s;
}
</style>