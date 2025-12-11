

export function fermitTime(time, k = false,y=0) {

  if (k == false) {
    var date = time.substr(0, 10); //年月日
    var hours = time.substring(11, 13);
    var minutes = time.substring(14, 16);
    var seconds = time.substring(17, 19);
    var timeFlag = date + " " + hours + ":" + minutes + ":" + seconds;
    timeFlag = timeFlag.replace(/-/g, "/");
    timeFlag = new Date(timeFlag);
    timeFlag = new Date(timeFlag.getTime() + 8 * 3600 * 1000);
    timeFlag =
      timeFlag.getFullYear() +
      "-" +
      (timeFlag.getMonth() + 1 < 10
        ? "0" + (timeFlag.getMonth() + 1)
        : timeFlag.getMonth() + 1) +
      "-" +
      (timeFlag.getDate() < 10 ? "0" + timeFlag.getDate() : timeFlag.getDate())
    return timeFlag;
  } else {
    var timeFlagt = new Date(time)
    var ss =
    (timeFlagt.getFullYear()+y) +
      "-" +
      (timeFlagt.getMonth() + 1 < 10
        ? "0" + (timeFlagt.getMonth() + 1)
        : timeFlagt.getMonth() + 1) +
      "-" +
      (timeFlagt.getDate() < 10 ? "0" + timeFlagt.getDate() : timeFlagt.getDate())
    return ss;
  }


}


export function getNextMonth(date) {
  var lastMonthToday = new Date(
    new Date().getTime()
  );
  var lastMonthYear = lastMonthToday.getFullYear() ;
  var lastMonth = lastMonthToday.getMonth();
  var lastMonthDay =
    lastMonthToday.getDate < 10
      ? "0" + lastMonthToday.getDate
      : lastMonthToday.getDate();
  var lastMonthKsrq = lastMonthYear + "-" + lastMonth + "-" + lastMonthDay;

  var nextMonthToday = new Date(
    new Date().getTime()
  );
  var nextMonthYear = nextMonthToday.getFullYear();
  var nextMonth = nextMonthToday.getMonth()+1;
  var nextMonthDay =
    nextMonthToday.getDate < 10
      ? "0" + nextMonthToday.getDate
      : nextMonthToday.getDate();
  var nextMonthJsrq = nextMonthYear + "-" + nextMonth + "-" + nextMonthDay;
  return nextMonthJsrq

}

//获取下个月的时间

export function getoneMonth(date) {
  let now = new Date();
  let year = now.getFullYear();
  let month = (now.getMonth() + 2);
  if (month >= 1 && month <= 9) {
    month = "0" + month;
  }
  if(month>12){
    month = month-12
    year = year+1
  }
  let firstDate = year + '年' +month + '月' + date + '日';
  return firstDate;

}

//获取开通天数
export const getDateByDays = (startDate, endDdate) => {
  let d1 = Date.parse(new Date(startDate));
  let d2 = Date.parse(new Date(endDdate));
  // 时间戳相减 / 天数
  let day = parseInt((d2 - d1) / (1000 * 60 * 60 * 24));
  return day + 1
}

//获取到期的时间

export function getoneyear(time, date) {
  let now = new Date(time);
  let year = now.getFullYear() + date;
  let month = now.getMonth() + 1;
  let day = now.getDate()
  if (month >= 1 && month <= 9) {
    month = "0" + month;
  }
  let firstDate = year + '年' + month + '月' + day + '日';
  return firstDate;
}

export function bjTime(time) {
  var chinaStandard = time
  var date = new Date(chinaStandard);
  var y = date.getFullYear();
  var m = date.getMonth() + 1;
  m = m < 10 ? ('0' + m) : m;
  var d = date.getDate();
  d = d < 10 ? ('0' + d) : d;
  var h = date.getHours();
  var minute = date.getMinutes();
  let second = date.getSeconds()
  minute = minute < 10 ? ('0' + minute) : minute;
  second = second < 10 ? ('0' + second) : second;
  let times = y + '-' + m + '-' + d + ' ' + h + ':' + minute + ':' + second;
  return times
}

export function getLastDayOfNaturalMonth(time) {
  // 获取当前日期
  const currentDate = new Date(time);
 
  // 获取当前年份
  let year = currentDate.getFullYear();
 
  // 获取下个月的月份，使用 padStart 在月份小于 10 时在前面补零
  let month = String(currentDate.getMonth() + 2).padStart(2, "0");
 
  // 如果月份为 '13'，说明当前是12月，需要更新年份并将月份设置为 '01'
  if (month === '13') {
    year = year + 1;
    month = '01';
  }
 
  // 设置每个月的第一天为 '01'
  let firstDay = "01";
 
  // 构建下个月的第一天的日期字符串，格式为 'YYYY-MM-DD'
  const firstDayOfNextMonth = `${year}-${month}-${firstDay}`;
 
  // 计算当前月的最后一天，通过减去一天的毫秒数来得到
  const lastDayOfMonth = new Date(new Date(firstDayOfNextMonth).getTime() - 86400000);
 
  // 将最后一天的日期转换为 ISO 格式，并提取日期部分
  return lastDayOfMonth
}
