# bento

弁当屋さんサイトからメニューを取得して Slack に通知します。
Google App Engine 用ウェブアプリです。

main.go で `slackHandler.url` に Slack の Incoming Webhook アプリから得た URL セットして、GAE にデプロイしてください。
 