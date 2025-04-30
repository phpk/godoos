<script setup lang="ts">
import {reactive, watch} from 'vue'
import {Solar, SolarMonth, SolarWeek, HolidayUtil} from 'lunar-typescript'

const now = Solar.fromDate(new Date())

class Day {
  public month: number = 0
  public day: number = 0
  public lunarDay: string = ''
  public lunarMonth: string = ''
  public yearGanZhi: string = ''
  public yearShengXiao: string = ''
  public monthGanZhi: string = ''
  public dayGanZhi: string = ''
  public ymd: string = ''
  public desc: string = ''
  public isToday: boolean = false
  public isSelected: boolean = false
  public isRest: boolean = false
  public isHoliday: boolean = false
  public festivals: string[] = []
  public yi: string[] = []
  public ji: string[] = []
}

class Week {
  public days: Day[] = []
}

class Month {
  public heads: string[] = []
  public weeks: Week[] = []
}

class Holiday {
  public name: string = ''
  public month: number = 0
}

const state = reactive({
  year: now.getYear(),
  month: now.getMonth(),
  weekStart: 1,
  selected: new Day(),
  data: new Month(),
  holidays: new Array<Holiday>(),
  holidayMonth: 0
})

function buildDay(d: Solar) {
  const ymd = d.toYmd()
  const lunar = d.getLunar()
  const day = new Day()
  day.month = d.getMonth()
  day.day = d.getDay()
  day.lunarMonth = lunar.getMonthInChinese()
  day.lunarDay = lunar.getDayInChinese()
  day.yearGanZhi = lunar.getYearInGanZhi()
  day.yearShengXiao = lunar.getYearShengXiao()
  day.monthGanZhi = lunar.getMonthInGanZhi()
  day.dayGanZhi = lunar.getDayInGanZhi()
  day.ymd = ymd
  day.isToday = ymd == now.toYmd()
  day.isSelected = ymd == state.selected.ymd
  if (day.isToday && state.selected.day === 0) {
    state.selected = day
  }
  const solarFestivals = d.getFestivals()
  solarFestivals.forEach(f => {
    day.festivals.push(f)
  })
  d.getOtherFestivals().forEach(f => {
    day.festivals.push(f)
  })
  lunar.getFestivals().forEach(f => {
    day.festivals.push(f)
  })
  lunar.getOtherFestivals().forEach(f => {
    day.festivals.push(f)
  })
  let rest = false
  if (d.getWeek() === 6 || d.getWeek() === 0) {
    rest = true
  }
  const holiday = HolidayUtil.getHoliday(ymd)
  if (holiday) {
    rest = !holiday.isWork()
  }
  day.isHoliday = !!holiday
  day.isRest = rest
  day.yi = lunar.getDayYi()
  day.ji = lunar.getDayJi()
  let desc = lunar.getDayInChinese()
  const jq = lunar.getJieQi()
  if (jq) {
    desc = jq
  } else if (lunar.getDay() === 1) {
    desc = lunar.getMonthInChinese() + '月'
  } else if (solarFestivals.length > 0) {
    const f = solarFestivals[0]
    if (f.length < 4) {
      desc = f
    }
  }
  day.desc = desc
  return day
}

function render() {
  const month = new Month()
  const weeks: SolarWeek[] = []
  const solarWeeks = SolarMonth.fromYm(parseInt(state.year + '', 10), parseInt(state.month + '', 10)).getWeeks(state.weekStart)
  solarWeeks.forEach(w => {
    weeks.push(w)
  })
  while (weeks.length < 6) {
    weeks.push(weeks[weeks.length - 1].next(1, false))
  }
  weeks.forEach(w => {
    const week = new Week()
    const heads: string[] = []
    w.getDays().forEach(d => {
      heads.push(d.getWeekInChinese())
      week.days.push(buildDay(d))
    })
    month.heads = heads
    month.weeks.push(week)
  })
  state.data = month
  const holidays: Holiday[] = []
  HolidayUtil.getHolidays(state.year).forEach(h => {
    const holiday = new Holiday()
    holiday.name = h.getName()
    holiday.month = parseInt(h.getTarget().substring(5, 7), 10)
    const exists = holidays.some(a => {
      return a.name == holiday.name
    })
    if (!exists) {
      holidays.push(holiday)
    }
  })
  state.holidays = holidays
}

function onSelect(day: Day) {
  state.selected = day
}

function onBack() {
  state.holidayMonth = 0
  state.year = now.getYear()
  state.month = now.getMonth()
  state.selected = buildDay(now)
}

render()

watch(() => state.year, () => {
  render()
})

watch(() => state.month, () => {
  render()
})

watch(() => state.selected, () => {
  render()
})

watch(() => state.holidayMonth, (newVal) => {
  const month = parseInt(newVal + '', 10)
  if (month > 0) {
    state.month = month
    render()
  }
})

</script>

<template>
  <div class="calendar">
    <div class="container">
      <div class="bar">
        <div>
          <input v-model="state.year">年
        </div>
        <div>
          <select v-model="state.month">
            <option :value="i" v-for="i in 12">{{ i }}月</option>
          </select>
        </div>
        <div>
          <select v-model="state.holidayMonth">
            <option value="0">假期安排</option>
            <option :value="h.month" v-for="h in state.holidays">{{ h.name }}</option>
          </select>
        </div>
        <div>
          <div class="button" @click="onBack">返回今天</div>
        </div>
      </div>
      <ul class="head">
        <li v-for="head in state.data.heads">{{ head }}</li>
      </ul>
      <ul class="body">
        <li v-for="week in state.data.weeks">
          <ol>
            <li @click="onSelect(day)" v-for="day in week.days"
                :class="{today: day.isToday, selected: day.isSelected, other: day.month != state.month, rest: day.isRest, holiday: day.isHoliday}">
              <div class="inner">
                <b>{{ day.day }}</b>
                <i>{{ day.desc }}</i>
                <u v-if="day.isHoliday">{{ day.isRest ? '休' : '班' }}</u>
              </div>
            </li>
          </ol>
        </li>
      </ul>
    </div>
    <div class="side">
      <div class="ymd">{{ state.selected.ymd }}</div>
      <div class="day">{{ state.selected.day }}</div>
      <div class="lunar">
        <div>{{ state.selected.lunarMonth }}月 {{ state.selected.lunarDay }}</div>
        <div>{{ state.selected.yearGanZhi }}年 {{ state.selected.yearShengXiao }}</div>
        <div>{{ state.selected.monthGanZhi }}月 {{ state.selected.dayGanZhi }}日</div>
      </div>
      <div class="festival" v-for="f in state.selected.festivals">
        {{ f }}
      </div>
      <div class="yiji">
        <div class="yi">
          <b>宜</b>
          <div v-for="f in state.selected.yi">
            {{ f }}
          </div>
        </div>
        <div class="ji">
          <b>忌</b>
          <div v-for="f in state.selected.ji">
            {{ f }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/calendar.scss"
</style>
