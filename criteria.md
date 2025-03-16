---
tags: []
creation_date: "2025-03-16"
---

## To Do:
- [x] Add test cases for the download functions, test if the download files have the correct full content.
- [x] Add list to catch error stocks in the download functions and a list to catch the stocks that download successfully.
- [x] Change combine csv function to only combine stocks that are in the success list.
- [ ] Add a function to add headers before combining, one for each sheet.
- [ ] Organize helper folder into more precise functions.
- [ ] Add test for each function file
    - [ ] helper
    - [ ] scraper
    - [ ] storage


## Sheet 1 Update
- Remove Columns G ~ L (Different PER Price)
- Add Index data to the same sheet. (Do it for each stock) (The dates should align with each other) (Leave one column in the middle blank)
    - Link: https://goodinfo.tw/tw/ShowK_Chart.asp?STOCK_ID=加權指數&CHT_CAT=WEEK&PRICE_ADJ=F&SCROLL2Y=0 
- After Index, leave one column blank and add the data from the following link for each stock. 
    - TSMC Link Example: https://goodinfo.tw/tw/ShowK_Chart.asp?STOCK_ID=2330&CHT_CAT=WEEK&PRICE_ADJ=F&SCROLL2Y=0

## Add Sheet 2
- Monthly Revenue.
- TSMC Example Link: https://goodinfo.tw/tw/ShowSaleMonChart.asp?STOCK_ID=2330
- Cash Flow Statement
- TSMC Example Link: https://goodinfo.tw/tw/StockCashFlow.asp?STOCK_ID=2330&PRICE_ADJ=F&SCROLL2Y=215&RPT_CAT=M_QUAR


## First Sheet
- PER on the left
    https://goodinfo.tw/tw/ShowK_ChartFlow.asp?RPT_CAT=PER&STOCK_ID=2330
- Stock Price on the right
    https://goodinfo.tw/tw/ShowK_Chart.asp?STOCK_ID=2330&CHT_CAT=WEEK&PRICE_ADJ=T&SCROLL2Y=389

## Second Sheet
- Monthly Revenue on the left
- Cash Flow Statement on the right
