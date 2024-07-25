// 初始化变量和常量
const icon_play = "&#xe769;";
const icon_stop = "&#xe637;";
const icon_study = "&#xe61e;";
const icon_rest = "&#xe623;";
let t_study = 45*60;
let t_rest = 5*60;
let plan_name = "godo发布";
let plan_time = "2021/12/31 18:00:00";
let isPlay = false;
let isRest = false;
let isShowPlan = false;
let isClock = false;
let localData = {};
let today = null;
let res_time = t_study;
let study_time = 0;
//获取现在的日期并返回
function getToday(){
    let t = new Date();
    let day = t.getDate();
    let month = t.getMonth() +1 ;
    let year = t.getFullYear();
    let tmp = year +
    (month < 10 ? "0" + month : month) +
    (day < 10 ? "0" + day : day);
    return tmp;
}
//更新today变量为最新的当前日期
function updateToday(){
    today = getToday();
}

//读取本地的存储，放到localDate对象中
function loadLocalData(){
    updateToday();
    //如果存在本地数据
    if(window.localStorage && window.localStorage.length){
        for (let i=0;i<window.localStorage.length;i++){
            let date = window.localStorage.key(i);
            let data = window.localStorage.getItem(date);
            localData[date] = data;
        }
        //console.log('本地数据读取完成：'+JSON.stringify(localData));
        //如果有今天的数据，那今天的工作时间叠加
        if(localData[today]){
            study_time+=localData[today];
        }
        if(localData['plan-name']){
            plan_name = localData['plan-name'];
            $(".p4 input").eq(0).val(plan_name);
        }
        if(localData['plan-time']){
            plan_time = localData['plan-time'];
            $(".p4 input").eq(1).val(plan_time.substr(0,10));
        }
        if(localData['t_rest']){
            t_rest = localData['t_rest'];
            $(".p3 input").eq(1).val(t_rest/60);
        }
        if(localData['t_study']){
            t_study = localData['t_study'];
            $(".p3 input").eq(0).val(t_study/60);
            res_time = t_study;
            showRestTime();
        }
    }
    else {
        study_time = 0;
    }
}
//把剩余的秒数计算成四个要显示的数字并返回
function calTime(t){
    let min = Math.floor(t/60);
    let sec = t %60;
    return [Math.floor(min/10) , min%10 , Math.floor(sec/10) , sec %10];
}
//更新计数器时间
function showRestTime(){
    let arr = calTime(res_time);
    $("#n1").text(arr[0]);
    if(arr[0]==1){
        $("#n1").css("transform","translate(30%)");
    }
    else {
        $("#n1").css("transform","none");
    }
    $("#n2").text(arr[1]);
    if(arr[1]==1){
        $("#n2").css("transform","translate(30%)");
    }
    else {
        $("#n2").css("transform","none");
    }
    $("#n4").text(arr[2]);
    if(arr[2]==1){
        $("#n4").css("transform","translate(30%)");
    }
    else {
        $("#n4").css("transform","none");
    }
    $("#n5").text(arr[3]);
    if(arr[3]==1){
        $("#n5").css("transform","translate(30%)");
    }
    else {
        $("#n5").css("transform","none");
    }
}
//更新现在的时间和电量
function showNowTime(){
    let t = new Date();
    let h = t.getHours();
    h = h < 10 ? "0" + h : h;
    let m = t.getMinutes();
    m = m < 10 ? "0" + m : m;
    let s = t.getSeconds();
    s = s < 10 ? "0" + s : s;
    if(navigator.getBattery){
        navigator.getBattery().then((result)=>{
            $(".now-time").text(h + ":" + m + ":" + s +"  电量:"+parseInt(result.level*100)+"%") ;
        })
    }
    else {
        $(".now-time").text(h + ":" + m + ":" + s);
    }
    //是否开启钟表模式
    if(isClock){
        //console.log(h+m);
        let arr = [h.toString()[0],h.toString()[1],m.toString()[0],m.toString()[1]]
        $("#n1").text(arr[0]);
        $("#n2").text(arr[1]);
        $("#n4").text(arr[2]);
        $("#n5").text(arr[3]);
        if(arr[0]==1){
            $("#n1").css("transform","translate(30%)");
        }
        else {
            $("#n1").css("transform","none");
        }
        $("#n2").text(arr[1]);
        if(arr[1]==1){
            $("#n2").css("transform","translate(30%)");
        }
        else {
            $("#n2").css("transform","none");
        }
        $("#n4").text(arr[2]);
        if(arr[2]==1){
            $("#n4").css("transform","translate(30%)");
        }
        else {
            $("#n4").css("transform","none");
        }
        $("#n5").text(arr[3]);
        if(arr[3]==1){
            $("#n5").css("transform","translate(30%)");
        }
        else {
            $("#n5").css("transform","none");
        }
    }
    else{
        showRestTime();
    }
}
//更新计划倒计时
function showPlan(){
    let t = new Date();
    let t0 = new Date(plan_time);
    let ss = t0.getTime() - t.getTime();
    let r_plan_days = parseInt(ss / 1000/60/60/24);
    let r_plan_hours = parseInt((ss / 1000/60/60)%24);
    $(".top-right").text("距" + plan_name+"还有:" + r_plan_days + "天" + r_plan_hours + "小时");
}
//存储工作时间
function saveStudyTime(){
    console.log('存储数据');
    if (getToday()!==today){
        study_time = 0 ;
        updateToday();
    }
    localData[today]=study_time.toFixed(2);
    window.localStorage.setItem(today,study_time);
}
//更新今天已工作时间
function showStudyTime(){
    $(".top-right").text("今天已工作:" +
    String(Math.floor(study_time / 60)) +
    "分钟" +
    String(study_time % 60) +
    "秒");
}

//更新最上面一栏
function showTopBar(){
    showNowTime();
    if(isShowPlan){
        showPlan();
    }
    else{
        showStudyTime();
    }
}
//保存数据
function saveSetting (){
    localData['plan-name'] = plan_name;
    localData['plan-time'] = plan_name;
    localData['t_study'] = t_study;
    localData['t_rest'] = t_rest;
    window.localStorage.setItem('plan-name',plan_name);
    window.localStorage.setItem('plan-time',plan_time);
    window.localStorage.setItem('t_study',t_study);
    window.localStorage.setItem('t_rest',t_rest);

}


//展示图表
function showEcharts(){
    //console.log("缓存数据:"+ window.localStorage);
    //console.log('localData'+localData);
    let myChart = echarts.init(document.getElementById('echarts'));
    let dateArr = [];
    let timeArr = [];
    let re = /^20[0-9]{2}/;
    for (let i in localData){
        if(re.test(i)){
            timeArr.push((localData[i]/60).toFixed(2));
            //处理原来的数据
            let re1 = /^0/;
            if (re1.test(i.substr(4,4))){
                dateArr.push(i[5]+'/'+i[6]+i[7]);
            }
            
        }
    }
    //处理数据补0
    if (dateArr.length<4){
        let n = 4 - dateArr.length;
        if (dateArr.length==0){
            let t = new Date();
            t = t.getMonth()+1 + '/'+t.getDate();
            dateArr.push(t);
            timeArr.push(0);
        }
        for (let j=0;j<n;j++){
            let t = new Date();
            t=t.setDate(t.getDate()+j+1);
            t = new Date(t);
            dateArr.push(t.getMonth()+1+'/'+t.getDate());
            timeArr.push(0);
        }
        
    }
    // console.log(dateArr);
    // console.log(timeArr);
    let option = {
        title: {
            text: '工作记录',
            subtext: '任意点击返回',
            subtextStyle: {
                align: 'right'
            },
            left: 'center',
            icon : 'none',
            textStyle: {
                fontSize: 26
            },
            
        },
        
    
        xAxis: {
            data: dateArr,
            name: '日期'
        },
        yAxis: {
            name: '工作时间/分钟',
            max: function(value){
                return value.max<30?30:value.max;
            }
        },
        series: [{
            name: '时间(分)',
            type: 'bar',
            data: timeArr,
            barMaxWidth: '40%',
            itemStyle:{
                color: '#555',
                normal: {
                    label: {
                        show: true,
                        position: 'top',
                        formatter: '{c}分钟'
                    }
                }
                
            }
            }]
    };
    myChart.setOption(option);
    return myChart;
}
