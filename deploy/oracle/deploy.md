# Oracle Cloud Always Free Deployment

## Setup Steps

### 1. Create Oracle Cloud Account
- Sign up at https://cloud.oracle.com/
- Get Always Free tier (no credit card required after trial)

### 2. Create Compute Instance
```bash
# Create VM.Standard.A1.Flex instance (ARM-based, Always Free)
# 1-4 OCPUs, 1-24 GB RAM
# 200 GB Block Storage
```

### 3. Install Dependencies
```bash
# Connect to your instance
ssh -i your-key.pem ubuntu@your-instance-ip

# Install Go
wget https://go.dev/dl/go1.22.5.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf go1.22.5.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install MongoDB (or use MongoDB Atlas)
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod
```

### 4. Deploy Application
```bash
# Clone your repository
git clone https://github.com/yourusername/asset-pulse.git
cd asset-pulse

# Build application
go build -o notify cmd/notify.go
go build -o database cmd/database.go

# Create environment file
cp .env.example .env
# Edit .env with your credentials

# Create systemd service for daily execution
sudo tee /etc/systemd/system/daily-tracker.service > /dev/null <<EOF
[Unit]
Description=Daily Investment Tracker
After=network.target

[Service]
Type=oneshot
User=ubuntu
WorkingDirectory=/home/ubuntu/asset-pulse
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin
ExecStart=/bin/bash -c 'source .env && ./notify now'
ExecStartPost=/bin/bash -c 'source .env && ./database save'

[Install]
WantedBy=multi-user.target
EOF

# Create timer for daily execution
sudo tee /etc/systemd/system/daily-tracker.timer > /dev/null <<EOF
[Unit]
Description=Run Daily Investment Tracker
Requires=daily-tracker.service

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
EOF

# Enable and start timer
sudo systemctl daemon-reload
sudo systemctl enable daily-tracker.timer
sudo systemctl start daily-tracker.timer

# Check status
sudo systemctl status daily-tracker.timer
```

### 5. Monitor
```bash
# View logs
journalctl -u daily-tracker.service -f

# Manual run
sudo systemctl start daily-tracker.service
```

## Cost: $0/month (Always Free)
- VM.Standard.A1.Flex: 1-4 OCPUs, 1-24 GB RAM
- 200 GB Block Storage
- Unlimited bandwidth (10 TB/month)
- Always Free - no expiration