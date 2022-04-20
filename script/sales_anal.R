library(googlesheets4)
library(glue)
library(lubridate)
library(dplyr)
library(knitr)
library(kableExtra)
# ---- manual config ------------------------------------------------------
# Please check next parameters are correct before running the script
prev.year = 2019  
gs.id = "1MqFAbr4zWftHKI2YU9vHPSUFeygfOWR621Akne7dAu8"

# ---- automatic config ---------------------------------------------------
key.date = last.sunday()
key.month = month(key.date)
week.from = day(key.date - 6)
week.to = day(key.date)
curr.year = year(key.date)
last.year = curr.year - 1

# ---- functions ---------------------------------------------------------
last.sunday <- function() {
  yoil <- weekdays(today(), abbreviate = TRUE)
  case_when(
    yoil == "월" ~ today() - 1,
    yoil == "화" ~ today() - 2,
    yoil == "수" ~ today() - 3,
    yoil == "목" ~ today() - 4,
    yoil == "금" ~ today() - 5,
    yoil == "토" ~ today() - 6,
    yoil == "일" ~ today() - 7,
  )
}

curr.week.sales <- function(curr, last, prev) {
  df.curr <- curr %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day >= week.from, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(curr.d = sum(D_Sales), curr.i = sum(I_Sales), curr.t = sum(G_Sales)) 
  
  df.last <- last %>% mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day >= week.from, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(last.d = sum(D_Sales), last.i = sum(I_Sales), last.t = sum(G_Sales))
  
  df.prev <- prev %>% mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day >= week.from, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(prev.d = sum(D_Sales), prev.i = sum(I_Sales), prev.t = sum(G_Sales))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = sprintf("%d.%d ~ %d.%d", key.month, week.from, key.month, week.to)) 
}

curr.month.sales <- function(curr, last, prev) {
  df.curr <- curr %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(curr.d = sum(D_Sales), curr.i = sum(I_Sales), curr.t = sum(G_Sales)) 
  
  df.last <- last %>% mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(last.d = sum(D_Sales), last.i = sum(I_Sales), last.t = sum(G_Sales))
  
  df.prev <- prev %>% mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(prev.d = sum(D_Sales), prev.i = sum(I_Sales), prev.t = sum(G_Sales))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = sprintf("%d.1 ~ %d.%d", key.month, key.month, week.to)) 
}

monthly.sales <- function(curr, last, prev) {
  df.curr <- curr %>% 
    filter(month(Date) < key.month) %>% 
    group_by(month) %>% 
    summarise(curr.d = sum(D_Sales), curr.i = sum(I_Sales), curr.t = sum(G_Sales))
  
  df.last <- last %>% 
    filter(month(Date) < key.month) %>% 
    group_by(month) %>% 
    summarise(last.d = sum(D_Sales), last.i = sum(I_Sales), last.t = sum(G_Sales))
  
  df.prev <- prev %>% 
    filter(month(Date) < key.month) %>% 
    group_by(month) %>% 
    summarise(prev.d = sum(D_Sales), prev.i = sum(I_Sales), prev.t = sum(G_Sales))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = as.character(month))
}

accum.sales <- function(curr, last, prev) {
  df.curr <- curr %>% 
    filter(Date <= key.date) %>% 
    summarise(to.mm = key.month, to.dd = week.to, curr.d = sum(D_Sales), curr.i = sum(I_Sales), curr.t = sum(G_Sales))
  
  df.last <- last %>% 
    filter(Date <= key.date) %>% 
    summarise(to.mm = key.month, to.dd = week.to, last.d = sum(D_Sales), last.i = sum(I_Sales), last.t = sum(G_Sales))
  
  df.prev <- prev %>% 
    filter(Date <= key.date) %>% 
    summarise(to.mm = key.month, to.dd = week.to, prev.d = sum(D_Sales), prev.i = sum(I_Sales), prev.t = sum(G_Sales))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = sprintf("1.1 ~ %d.%d", to.mm, to.dd)) %>% 
    select(month, 3:11)
}

curr.week.sales.yr <- function(curr, last, prev) {
  df.curr <- curr %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day >= week.from, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(curr.d = sum(D_Sales + D_YR), 
              curr.i = sum(I_Sales + I_YR), 
              curr.t = sum(D_Sales + D_YR + I_Sales + I_YR)) 
  
  df.last <- last %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day >= week.from, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(last.d = sum(D_Sales + D_YR), 
              last.i = sum(I_Sales + I_YR), 
              last.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.prev <- prev %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day >= week.from, day <= week.to) %>%
    group_by(month) %>% 
    summarise(prev.d = sum(D_Sales + D_YR), 
              prev.i = sum(I_Sales + I_YR), 
              prev.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = sprintf("%d.%d ~ %d.%d", key.month, week.from, key.month, week.to)) 
}

curr.month.sales.yr <- function(curr, last, prev) {
  df.curr <- curr %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(curr.d = sum(D_Sales + D_YR), 
              curr.i = sum(I_Sales + I_YR), 
              curr.t = sum(D_Sales + D_YR + I_Sales + I_YR)) 
  
  df.last <- last %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day <= week.to) %>% 
    group_by(month) %>% 
    summarise(last.d = sum(D_Sales + D_YR), 
              last.i = sum(I_Sales + I_YR), 
              last.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.prev <- prev %>% 
    mutate(month = month(Date), day = day(Date)) %>%
    filter(month == key.month, day <= week.to) %>%
    group_by(month) %>% 
    summarise(prev.d = sum(D_Sales + D_YR), 
              prev.i = sum(I_Sales + I_YR), 
              prev.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = sprintf("%d.%1 ~ %d.%d", key.month, key.month, week.to)) 
}

monthly.sales.yr <- function(curr, last, prev) {
  df.curr <- curr %>% 
    filter(month(Date) < key.month) %>%
    group_by(month) %>% 
    summarise(curr.d = sum(D_Sales + D_YR), 
              curr.i = sum(I_Sales + I_YR), 
              curr.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.last <- last %>% 
    filter(month(Date) < key.month) %>% 
    group_by(month) %>% 
    summarise(last.d = sum(D_Sales + D_YR), 
              last.i = sum(I_Sales + I_YR), 
              last.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.prev <- prev %>% 
    filter(month(Date) < key.month) %>% 
    group_by(month) %>% 
    summarise(prev.d = sum(D_Sales + D_YR), 
              prev.i = sum(I_Sales + I_YR), 
              prev.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = as.character(month))
}

accum.sales.yr <- function(curr, last, prev) {
  df.curr <- curr %>% 
    filter(Date <= key.date) %>% 
    summarise(to.mm = key.month, to.dd = week.to, 
              curr.d = sum(D_Sales + D_YR), 
              curr.i = sum(I_Sales + I_YR), 
              curr.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.last <- last %>% 
    filter(Date <= key.date) %>% 
    summarise(to.mm = key.month, to.dd = week.to, 
              last.d = sum(D_Sales + D_YR), 
              last.i = sum(I_Sales + I_YR), 
              last.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.prev <- prev %>% 
    filter(Date <= key.date) %>% 
    summarise(to.mm = key.month, to.dd = week.to, 
              prev.d = sum(D_Sales + D_YR), 
              prev.i = sum(I_Sales + I_YR), 
              prev.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.curr %>% left_join(df.last) %>% left_join(df.prev) %>%
    mutate(month = sprintf("1.1 ~ %d.%d", to.mm, to.dd)) %>% 
    select(month, 3:11)
}

make.sales.smry <- function(curr, last, prev) {
  df.curr.week <- curr.week.sales(sales.curr, sales.last, sales.prev)
  df.monthly <- monthly.sales(sales.curr, sales.last, sales.prev)
  df.curr.month <- curr.month.sales(sales.curr, sales.last, sales.prev)
  df.accum   <- accum.sales(sales.curr, sales.last, sales.prev)
  bind_rows(df.curr.week, df.monthly, df.curr.month, df.accum)
}

make.sales.yr.smry <- function(curr, last, prev) {
  df.curr.week <- curr.week.sales.yr(sales.curr, sales.last, sales.prev)
  df.monthly <- monthly.sales.yr(sales.curr, sales.last, sales.prev)
  df.curr.month <- curr.month.sales.yr(sales.curr, sales.last, sales.prev)
  df.accum   <- accum.sales.yr(sales.curr, sales.last, sales.prev)
  bind_rows(df.curr.week, df.monthly, df.curr.month, df.accum)
}

# ---- sales results ----------------------------------------------------
sales.curr <- read_sheet(gs.id, sheet = glue('{curr.year}')) %>% 
  mutate(Date = ymd(format(Date, "%Y-%m-%d"))) %>% 
  mutate(month = month(Date), day = day(Date))

sales.last <- read_sheet(gs.id, sheet = glue('{last.year}')) %>% 
  mutate(Date = ymd(format(Date, "%Y-%m-%d"))) %>% 
  mutate(month = month(Date), day = day(Date))

sales.prev <- read_sheet(gs.id, sheet = glue('{prev.year}')) %>% 
  mutate(Date = ymd(format(Date, "%Y-%m-%d"))) %>% 
  mutate(month = month(Date), day = day(Date))

# n/a -> 0
sales.curr[is.na(sales.curr)] = 0
sales.last[is.na(sales.last)] = 0
sales.prev[is.na(sales.prev)] = 0

smry.sales <- make.sales.smry(sales.curr, sales.last, sales.prev)
smry.sales.yr <- make.sales.yr.smry(sales.curr, sales.last, sales.prev)


smry.sales %>% 
  mutate(curr.d = curr.d / 100000000, 
         curr.i = curr.i / 100000000,
         curr.t = curr.t / 100000000, last.d = last.d / 100000000,
         last.i = last.i / 100000000, last.t = last.t / 100000000,
         prev.d = prev.d / 100000000, prev.i = prev.i / 100000000,
         prev.t = prev.t / 100000000,
         d.YoY1 = sprintf("%+.1f%%", (curr.d / last.d - 1) * 100),
         d.YoY2 = sprintf("%+.1f%%", (curr.d / prev.d - 1) * 100),
         i.YoY1 = sprintf("%+.1f%%", (curr.i / last.i - 1) * 100),
         i.YoY2 = sprintf("%+.1f%%", (curr.i / prev.i - 1) * 100),
         t.YoY1 = sprintf("%+.1f%%", (curr.t / last.t - 1) * 100),
         t.YoY2 = sprintf("%+.1f%%", (curr.t / prev.t - 1) * 100)) %>%
  select(month, curr.d, d.YoY1, d.YoY2, curr.i, i.YoY1, i.YoY2, curr.t, t.YoY1, t.YoY2,
         last.d, last.i, last.t, prev.d, prev.i, prev.t)
  





