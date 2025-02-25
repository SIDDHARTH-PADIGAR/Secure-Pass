# 🔐 SecurePass – A Simple Yet Powerful Password Manager  

SecurePass is a lightweight and secure password manager designed to store, manage, and retrieve passwords efficiently. It provides a **command-line interface** for seamless password management while ensuring strong encryption and security.  

## ✨ Features  
- **🔑 Secure Password Storage** – Encrypts and stores passwords securely.  
- **📜 Password History Tracking** – Keeps a history of password changes for better account management.  
- **🔍 Search Functionality** – Quickly find stored passwords using keywords.  
- **📅 Expiry Alerts** – Notifies when a password is about to expire.  
- **📋 Clipboard Integration** – Copy passwords securely to the clipboard with auto-clear after 30 seconds.  
- **🔢 Strong Password Generator** – Generate complex passwords with customizable rules.  

## 🛠️ Tech Stack  
- **Go** – Core programming language.  
- **BuntDB** – Lightweight key-value database for storing encrypted passwords.  
- **x/crypto** – Secure encryption for password storage.  
- **Atotto Clipboard** – Handles secure clipboard operations.  

## 🚀 Getting Started  
### 1️⃣ Clone the Repository  
```sh
git clone https://github.com/your-username/secure-pass.git
cd secure-pass
```

### 2️⃣ Install Dependencies  
Ensure you have Go installed (v1.23+). Then, run:  
```sh
go mod tidy
```

### 3️⃣ Build and Run  
```sh
go run main.go
```

## 📌 Future Plans  
- **🖥 GUI Integration** using Fyne for a user-friendly experience.  
- **🌐 Cloud Sync** to securely store passwords across devices.  
- **🔓 Biometric Authentication** for enhanced security.
- **🛡️ Breach Monitoring** Checks if stored passwords have been compromised in known data breaches.  
