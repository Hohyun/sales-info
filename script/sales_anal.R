library(googlesheets4)
library(glue)
library(lubridate)
library(dplyr)
library(knitr)
library(kableExtra)

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

# ---- manual config ------------------------------------------------------
# Please check next parameters are correct before running the script
prev.year = 2019  
gs.source.id = "1MqFAbr4zWftHKI2YU9vHPSUFeygfOWR621Akne7dAu8"
gs.report.id = "1sn2miEacsxPRdaqmyMSuKQNOON_DVRp5of3LI0vf5d8"

# ---- automatic config ---------------------------------------------------
key.date = last.sunday() + 2
# key.date = ymd("20220419")
key.month = month(key.date)
week.from = day(key.date - 6)
week.to = day(key.date)
curr.year = year(key.date)
last.year = curr.year - 1

# ---- functions ---------------------------------------------------------
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
    filter(Date <= make_date(last.year, key.month, week.to)) %>% 
    summarise(to.mm = key.month, to.dd = week.to, last.d = sum(D_Sales), last.i = sum(I_Sales), last.t = sum(G_Sales))
  
  df.prev <- prev %>% 
    filter(Date <= make_date(prev.year, key.month, week.to)) %>% 
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
    mutate(month = sprintf("%d.1 ~ %d.%d", key.month, key.month, week.to)) 
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
    filter(Date <= make_date(last.year, key.month, week.to)) %>% 
    summarise(to.mm = key.month, to.dd = week.to, 
              last.d = sum(D_Sales + D_YR), 
              last.i = sum(I_Sales + I_YR), 
              last.t = sum(D_Sales + D_YR + I_Sales + I_YR))
  
  df.prev <- prev %>% 
    filter(Date <= make_date(prev.year, key.month, week.to)) %>% 
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
  df.accum   <- accum.sales(sales.curr, sales.last, sales.prev)
  df.monthly <- monthly.sales(sales.curr, sales.last, sales.prev)
  df.curr.month <- curr.month.sales(sales.curr, sales.last, sales.prev)
  bind_rows(df.curr.week, df.accum, df.monthly, df.curr.month)
}

make.sales.yr.smry <- function(curr, last, prev) {
  df.curr.week <- curr.week.sales.yr(sales.curr, sales.last, sales.prev)
  df.accum   <- accum.sales.yr(sales.curr, sales.last, sales.prev)
  df.monthly <- monthly.sales.yr(sales.curr, sales.last, sales.prev)
  df.curr.month <- curr.month.sales.yr(sales.curr, sales.last, sales.prev)
  bind_rows(df.curr.week, df.accum, df.monthly, df.curr.month)
}

# ---- sales results ----------------------------------------------------
sales.curr <- read_sheet(gs.source.id, sheet = glue('{curr.year}'),
                         col_types = "D_n_____nnn___n")  %>% 
  mutate(month = month(Date), day = day(Date))

sales.last <- read_sheet(gs.source.id, sheet = glue('{last.year}'),
                         col_types = "D_n_____nnn___n")  %>% 
  mutate(month = month(Date), day = day(Date))

sales.prev <- read_sheet(gs.source.id, sheet = glue('{prev.year}'),
                         col_types = "D_n_____nnn___n")  %>% 
  mutate(month = month(Date), day = day(Date))

# n/a -> 0
sales.curr[is.na(sales.curr)] = 0
sales.last[is.na(sales.last)] = 0
sales.prev[is.na(sales.prev)] = 0

smry.sales <- make.sales.smry(sales.curr, sales.last, sales.prev)
smry.sales.yr <- make.sales.yr.smry(sales.curr, sales.last, sales.prev)


write_sheet(smry.sales, gs.report.id, sheet = "salesdata_no_fsc")
write_sheet(smry.sales.yr, gs.report.id, sheet = "salesdata_with_fsc")





