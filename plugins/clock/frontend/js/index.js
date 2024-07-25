

// 入口函数
$(()=>{
    let myChart ;
    let isFullScreen = false;
    //屏幕常量！
    let noSleep = new NoSleep();
    let playSound = (src) => {
        au = document.createElement("audio");
 	    au.preload="auto";
        au.src = '/admin/tools/clock/audio/' + src + '.mp3';
        au.play();
    }
    //绑定事件
    // 全屏
    $(".main").on("click",()=>{
        if(isFullScreen){
            document.exitFullscreen();
        }
        else{
            document.body.requestFullscreen();
        }
        isFullScreen = !isFullScreen;
    });
    //绑定提示取消
    $(".p0").on("click",()=>{
        $(".page").hide();
        $(".p1").show();
        noSleep.enable();
    });
    $(".top-left").on("click",()=>{
        $(".page").hide();
        $(".p0").show();
    });
    //切换工作时间/计划倒计时
    $(".top-right").on("click",()=>{
        isShowPlan = !isShowPlan;
    });
    //切换时钟
    $(".top-center").on("click",()=>{
        isClock = !isClock;
        showNowTime();
    });
    //播放按钮
    $(".play").on("click",()=>{
        //如果没播放并且在休息，按两下直接开始工作
        if (!isPlay && isRest){
            res_time = t_study;
            isRest = false;
        }
        if (!isPlay){
            $(".play").html(icon_stop);
            playSound('tick');
        }
        if (isPlay){
            $(".play").html(icon_play);
        }
        
        isPlay = !isPlay;
    });
    //重置时间
    $(".reset").on("click",()=>{
        if (isRest){
            res_time = t_rest ;
        }
        else {
            res_time = t_study ;
        }
        showRestTime();
    });
    //设置
    $(".setting").on("click",()=>{
        $(".page").hide();
        $(".p2").show();
    });
    //设置工作/休息时间
    $(".p2 .option").eq(0).on("click",()=>{
        $(".page").hide();
        $(".p3").show();
    });
    //设置计划倒计时
    $(".p2 .option").eq(1).on("click",()=>{
        $(".page").hide();
        $(".p4").show();
    });
    //查看工作历史
    $(".p2 .option").eq(2).on("click",()=>{
        $(".page").hide();
        $(".p5").show();
        myChart = showEcharts();
    });
    //清除本地数据
    $(".p2 .option").eq(3).on("click",()=>{
        window.localStorage.clear();
        localData = [];
        

        alert("已清除！");
    });
    //返回
    $(".p2 .option").eq(4).on("click",()=>{
        $(".page").hide();
        $(".p1").show();
    });
    //第三页设置工作/休息时间
    $(".p3-bottom button").eq(0).on("click",()=>{
        t_study = $(".p3 input").eq(0).val() * 60;
        t_rest = $(".p3 input").eq(1).val() * 60;
        if(isRest){
            res_time = t_rest;
        }
        else {
            res_time = t_study;
        }
        showRestTime();
        saveSetting();
        // alert("完成!");
        // console.log(t_study);
    });
    $(".p3-bottom button").eq(1).on("click",()=>{
        $(".page").hide();
        $(".p2").show();
    });
    //第四页 设置计划倒计时
    $(".p4-bottom button").eq(0).on("click",()=>{
        plan_name = $(".p4 input").eq(0).val();
        plan_time = $(".p4 input").eq(1).val() + " " + "08:00:00";
        if(isShowPlan){
            showPlan();
        }
        saveSetting();
        // alert("完成!");
        // console.log(plan_name);
    });
    $(".p4-bottom button").eq(1).on("click",()=>{
        $(".page").hide();
        $(".p2").show();
    });
    //第五页返回
    $(".p5").on("click",()=>{
        $(".page").hide();
        $(".p2").show();
        myChart.clear();
    });
 
    

    //开始初始化,设置计时器
    
    loadLocalData();
    
    setInterval(()=>{
        showTopBar();
        //如果倒计时归零
        if(!res_time){
            isRest = !isRest;
            //该休息了
            if(isRest){
                $(".top-left").html(icon_rest+"休息中");
                playSound('rest');
                res_time = t_rest;
            }
            //该工作了
            else{
                $(".top-left").html(icon_study+"工作中");
                playSound('work');
                res_time = t_study;
            }
        }
        //如果在播放
        if (isPlay && !isClock){
            if(!isRest){
                $(".top-left").html(icon_study+"工作中");
                study_time++;
            }
            res_time--;
            saveStudyTime();
            showRestTime(); 
        }
    },1000);
});