# app-go
```bash
curl -X POST "https://api.telegram.org/bot<TOKENBOT>/setWebhook?url=https://crecs-bot.it/telegram&secret_token=INSERISCI_WEBHOOK_SECRET"
```

```bash
sudo kubectl -n crecs create secret generic telegram-secrets   --from-literal=botToken="<TOKENBOT>"   --from-literal=webhookSecret="INSERISCI_WEBHOOK_SECRET"
```