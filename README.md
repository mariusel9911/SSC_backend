# Aplicație de Autentificare cu Factori Multipli în Go

## Tema Generală

Această aplicație implementează un sistem de autentificare cu factori multipli (Multi-Factor Authentication - MFA) folosind limbajul de programare Go. Sistemul este compus dintr-un frontend și un backend care interacționează pentru a oferi o soluție completă de autentificare securizată, utilizând multiple metode de verificare pentru a confirma identitatea utilizatorilor.

## Adresele de Repository

- Frontend: [https://github.com/mariusel9911/SSC_frontend](https://github.com/mariusel9911/SSC_frontend)
- Backend: [https://github.com/mariusel9911/SSC_backend](https://github.com/mariusel9911/SSC_backend)

## Documentație

### Concepte Generale

#### Autentificarea cu Factori Multipli (MFA)

Autentificarea cu factori multipli (MFA) reprezintă o metodă de securitate care necesită mai multe forme de verificare înainte de a acorda accesul unui utilizator. Acest sistem îmbunătățește semnificativ securitatea aplicațiilor față de metoda tradițională de autentificare bazată doar pe utilizator și parolă.

În general, factorii de autentificare pot fi clasificați în trei categorii principale:
1. **Ceva ce știi** - informații precum parole sau răspunsuri la întrebări de securitate
2. **Ceva ce ai** - dispozitive fizice precum un telefon mobil sau un token de securitate
3. **Ceva ce ești** - date biometrice precum amprente, recunoaștere facială sau scanare retiniană

Aplicația noastră implementează primele două categorii, utilizând parola ca primul factor și un cod temporar generat prin metoda TOTP (Time-based One-Time Password) ca al doilea factor.

#### TOTP (Time-based One-Time Password)

TOTP este o metodă standardizată (RFC 6238) pentru generarea de coduri de unică folosință bazate pe timp. Principalele caracteristici ale TOTP includ:

- Generarea codurilor pe baza unui secret partajat între server și dispozitivul utilizatorului
- Codurile sunt valabile pentru o perioadă limitată de timp (de obicei 30 de secunde)
- Algoritmul folosește timpul curent ca input, asigurând generarea de coduri diferite la fiecare interval de timp
- Implementarea este compatibilă cu aplicații populare precum Google Authenticator, Authy și Microsoft Authenticator

#### Sesiuni Web

Pentru gestionarea autentificării, aplicația utilizează sesiuni web tradiționale:

- Sesiunile sunt stocate server-side și identificate prin cookie-uri
- Cookie-urile conțin doar ID-ul sesiunii, nu informații sensibile
- Sesiunile permit invalidarea imediată la logout
- Implementarea este mai simplă și directă pentru aplicații web tradiționale
- Sesiunile sunt configurate cu parametri de securitate (HttpOnly, SameSite)

### Arhitectură și Design

#### Arhitectura Generală

Aplicația este structurată folosind un model client-server, cu separare clară între frontend și backend:

1. **Frontend** - interfața cu utilizatorul, dezvoltată cu React și găzduită în repository-ul SSC_frontend
2. **Backend** - serviciul API care gestionează logica de autentificare, dezvoltat în Go și găzduit în repository-ul SSC_backend

Comunicarea între cele două componente se realizează prin API RESTful cu autentificare pe bază de sesiuni.

#### Principii arhitecturale

Aplicația implementează următoarele principii de arhitectură:

1. **Separation of Concerns (SoC)** - Separarea clară a responsabilităților între componente, cu frontend-ul responsabil pentru interfața utilizator și backend-ul pentru logica de business și securitate.

2. **Backend for Frontend (BFF)** - Un model arhitectural în care backend-ul este special conceput pentru a servi nevoile frontend-ului, optimizând comunicarea între cele două componente.

3. **Arhitectură pe straturi** - Backend-ul este structurat în straturi logice (prezentare, servicii, date) pentru a separa responsabilitățile și a facilita testarea și mentenanța.

4. **Model API-First** - Dezvoltarea a început cu definirea clară a API-ului, care a servit ca un contract între frontend și backend.

5. **Session-based Authentication** - Utilizarea sesiunilor web pentru a implementa autentificarea, oferind o experiență tradițională și sigură pentru utilizatori.

#### Diagrama arhitecturală

```
+----------------+      HTTPS       +----------------+
|                |  <-------------> |                |
|    Frontend    |                  |    Backend     |
|   (Browser)    |  REST API calls  |   (Go Server)  |
|                |  <-------------> |                |
+----------------+                  +-------+--------+
                                            |
                                            | 
                                    +-------v--------+
                                    |                |
                                    |  Database      |
                                    |  (PostgreSQL)  |
                                    |                |
                                    +----------------+
```

#### Componente Backend

Backend-ul este construit în Go și include următoarele componente principale:

1. **Server HTTP** - gestionează rutele API și cererile HTTP, implementat folosind biblioteca Gin Gonic pentru crearea unui server HTTP performant și ușor de utilizat.

2. **Controller de Autentificare** - procesează cererile de înregistrare, autentificare și validare MFA:
   - `AuthController`: Gestionează înregistrarea și autentificarea utilizatorilor
   - `OTPController`: Gestionează generarea, validarea și dezactivarea tokenurilor TOTP

3. **Middleware** - Componente intermediare pentru procesarea cererilor:
   - `Session Management`: Gestionează sesiunile utilizatorilor
   - `CORS Configuration`: Configurează politicile de cross-origin resource sharing

4. **Generator TOTP** - creează și validează codurile temporare folosind algoritmi standardizați:
   - Implementat folosind biblioteca `pquerna/otp` care respectă standardul RFC 6238
   - Generează secrete unice pentru fiecare utilizator
   - Validează codurile cu toleranță temporală pentru a compensa decalajele de timp între dispozitive

5. **Servicii** - Implementează logica de business:
   - `UserService`: Gestionează operațiunile legate de utilizatori
   - `SessionService`: Gestionează sesiunile utilizatorilor
   - `OTPService`: Gestionează operațiunile legate de autentificarea cu doi factori

6. **Stocare date** - gestionează persistența informațiilor:
   - Utilizează o bază de date PostgreSQL pentru stocarea datelor utilizatorilor
   - Implementează operațiuni CRUD pentru utilizatori și setări MFA
   - Stochează secretele TOTP în format securizat

7. **Utilitare**:
   - Logare structurată pentru monitorizarea activității
   - Gestionare configurație din variabile de mediu
   - Module helper pentru validarea datelor

Biblioteca principală utilizată pentru implementarea TOTP este bazată pe standardul RFC 6238, care definește algoritmi OATH-TOTP pentru generarea codurilor temporare.

#### Structura directorului backend

```
SSC_backend/
├── main.go                         # Punctul de intrare al aplicației
├── config/
│   └── config.go                   # Configurația aplicației și conexiunea la DB
├── controllers/
│   └── auth.go                     # Controllere pentru autentificare și TOTP
├── models/
│   └── user.go                     # Modelul de date pentru utilizatori
├── routes/
│   └── routes.go                   # Definirea rutelor API
├── scripts/
│   └── resetdb.sh                  # Script pentru resetarea bazei de date
├── .env.example                    # Exemplu de configurare variabile de mediu
├── docker-compose.yml              # Configurare Docker pentru dezvoltare
├── Dockerfile                      # Configurare containerizare aplicație
├── go.mod                          # Dependențe Go
└── go.sum                          # Verificare integritate dependențe
```

#### Componente Frontend

Frontend-ul este dezvoltat utilizând React și TypeScript pentru o experiență de utilizare modernă și sigură. Principalele componente includ:

1. **Framework UI** - Implementat folosind tehnologii moderne de frontend:
   - Framework React cu TypeScript pentru type safety
   - Zustand pentru state management
   - React Router pentru navigarea client-side
   - Tailwind CSS pentru styling

2. **Module de autentificare**:
   - `AuthProvider`: Componentă React care gestionează starea de autentificare
   - `useStore`: Hook pentru accesul la starea globală a aplicației
   - `SessionManager`: Gestionează comunicarea cu backend-ul pentru sesiuni

3. **Componente de interfață**:
   - `RegisterForm`: Formular pentru înregistrarea utilizatorilor noi
   - `LoginForm`: Formular pentru autentificarea cu utilizator și parolă
   - `TwoFactorAuth`: Modal pentru configurarea autentificării cu doi factori
   - `OTPValidator`: Componentă pentru introducerea și validarea codurilor TOTP
   - `ProfilePage`: Interfață pentru gestionarea profilului utilizatorului

4. **Servicii**:
   - `authApi`: Configurează axios pentru comunicarea cu backend-ul
   - `QRCode`: Generează coduri QR pentru scanarea cu aplicații de autentificare
   - `ToastNotifications`: Afișează notificări și mesaje pentru utilizator

5. **Utilitare**:
   - Module pentru validarea formularelor cu Zod
   - Gestionarea erorilor și feedback-ului utilizatorului
   - Componente reutilizabile (FormInput, LoadingButton, Spinner)

Frontend-ul oferă interfețe complete pentru:

1. **Înregistrare utilizator** - cu validare în timp real și feedback vizual
2. **Autentificare primară** (cu utilizator și parolă)
3. **Configurare MFA** - inclusiv generarea și scanarea codurilor QR pentru aplicații de autentificare
4. **Validare MFA** - introducerea codului temporar pentru finalizarea autentificării
5. **Gestionare profil** - inclusiv activarea/dezactivarea MFA

#### Structura directorului frontend

```
SSC_frontend/
├── public/
│   ├── index.html                  # Template HTML principal
│   └── vite.svg                    # Icon aplicație
├── src/
│   ├── components/
│   │   ├── FormInput.tsx           # Componentă input cu validare
│   │   ├── Header.tsx              # Header cu navigare
│   │   ├── Layout.tsx              # Layout principal
│   │   ├── LoadingButton.tsx       # Buton cu stare de loading
│   │   ├── Spinner.tsx             # Componentă loading
│   │   └── TwoFactorAuth.tsx       # Modal pentru configurare 2FA
│   ├── pages/
│   │   ├── home.page.tsx           # Pagina principală
│   │   ├── login.page.tsx          # Pagina de autentificare
│   │   ├── profile.page.tsx        # Pagina de profil
│   │   ├── register.page.tsx       # Pagina de înregistrare
│   │   └── validate2fa.page.tsx    # Pagina de validare TOTP
│   ├── api/
│   │   ├── authApi.ts              # Configurare axios
│   │   └── types.ts                # Tipuri TypeScript
│   ├── router/
│   │   └── index.tsx               # Configurare routing
│   ├── store/
│   │   └── index.ts                # State management cu Zustand
│   ├── utils/
│   │   └── errorHandler.ts         # Gestionarea erorilor
│   ├── App.tsx                     # Componenta principală
│   ├── main.tsx                    # Punctul de intrare
│   └── index.css                   # Stiluri globale
├── package.json                    # Dependențe și scripturi
├── tailwind.config.cjs             # Configurare Tailwind CSS
├── tsconfig.json                   # Configurare TypeScript
└── vite.config.ts                  # Configurare Vite
```

### Fluxul de Autentificare

#### Diagrama fluxului de autentificare

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │     │             │
│  Utilizator │     │  Frontend   │     │   Backend   │     │   Bază de   │
│             │     │             │     │             │     │    date     │
│             │     │             │     │             │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │                   │
       │  Acces aplicație  │                   │                   │
       │──────────────────>│                   │                   │
       │                   │                   │                   │
       │                   │   Verifică sesiune│                   │
       │                   │───────────────────>                   │
       │                   │                   │                   │
       │                   │  Sesiune invalidă │                   │
       │                   │<───────────────────                   │
       │   Redirecționare  │                   │                   │
       │     spre login    │                   │                   │
       │<──────────────────│                   │                   │
       │                   │                   │                   │
       │  Introduce credențiale                │                   │
       │──────────────────>│ Cerere autentificare                  │
       │                   │───────────────────>                   │
       │                   │                   │                   │
       │                   │                   │  Verifică credențiale
       │                   │                   │──────────────────>│
       │                   │                   │                   │
       │                   │                   │ Credențiale valide│
       │                   │                   │<──────────────────│
       │                   │                   │                   │
       │                   │                   │  Verifică MFA activat
       │                   │                   │──────────────────>│
       │                   │                   │                   │
       │                   │                   │  MFA este activat │
       │                   │                   │<──────────────────│
       │                   │                   │                   │
       │                   │ Solicită cod OTP  │                   │
       │                   │<───────────────────                   │
       │                   │                   │                   │
       │ Solicită cod OTP  │                   │                   │
       │<──────────────────│                   │                   │
       │                   │                   │                   │
       │ Introduce cod OTP │                   │                   │
       │──────────────────>│  Validare cod OTP │                   │
       │                   │───────────────────>                   │
       │                   │                   │  Verifică secret  │
       │                   │                   │──────────────────>│
       │                   │                   │                   │
       │                   │                   │  Obține secret    │
       │                   │                   │<──────────────────│
       │                   │                   │                   │
       │                   │                   │  Validează OTP    │
       │                   │                   │  cu secret        │
       │                   │                   │                   │
       │                   │   OTP valid       │                   │
       │                   │<───────────────────                   │
       │                   │                   │                   │
       │                   │  Creează sesiune  │                   │
       │                   │<───────────────────                   │
       │                   │                   │                   │
       │  Acces acordat    │                   │                   │
       │<──────────────────│                   │                   │
       │                   │                   │                   │
```

#### Procese detaliate

1. **Înregistrare**:
   - Utilizatorul accesează pagina de înregistrare și completează formularul cu datele personale (nume, email, parolă)
   - Frontend-ul validează datele introduse folosind Zod pentru validare și React Hook Form pentru gestionarea formularelor
   - Datele validate sunt trimise către backend prin API
   - Backend-ul verifică unicitatea email-ului în baza de date
   - Parola este hash-uită folosind bcrypt cu cost factor adecvat pentru securitate
   - Datele utilizatorului sunt stocate în baza de date PostgreSQL
   - Utilizatorul este redirecționat către pagina de login pentru autentificare

2. **Activare MFA**:
   - Utilizatorul autentificat accesează pagina de profil și selectează opțiunea de activare MFA
   - Frontend-ul trimite o cerere către backend pentru inițierea procesului de configurare
   - Backend-ul generează un secret TOTP unic pentru utilizator folosind biblioteca pquerna/otp
   - Secretul este salvat în baza de date în câmpul OTPSecret al utilizatorului
   - Backend-ul creează un URL otpauth:// compatibil cu aplicațiile de autentificare
   - Frontend-ul primește URL-ul și secretul, generând un cod QR folosind biblioteca qrcode
   - Utilizatorul scanează codul QR cu o aplicație de autentificare (Google Authenticator, Authy, etc.)
   - Aplicația de autentificare stochează secretul și începe să genereze coduri TOTP
   - Utilizatorul introduce codul curent afișat în aplicație pentru verificare
   - Backend-ul validează codul folosind totp.Validate() și activează MFA pentru utilizator
   - Sesiunea utilizatorului este actualizată pentru a reflecta activarea MFA

3. **Autentificare cu MFA**:
   - Utilizatorul accesează pagina de autentificare și introduce email și parolă (primul factor)
   - Frontend-ul trimite credențialele către backend prin API
   - Backend-ul validează credențialele folosind bcrypt.CompareHashAndPassword()
   - Backend-ul verifică dacă utilizatorul are MFA activat (OTPEnabled = true)
   - Pentru utilizatorii cu MFA activat, backend-ul răspunde cu indicator otp:true
   - Frontend-ul redirecționează utilizatorul către pagina de introducere cod TOTP
   - Utilizatorul consultă aplicația de autentificare și introduce codul generat de 6 cifre
   - Frontend-ul trimite codul către backend pentru validare
   - Backend-ul recuperează secretul TOTP al utilizatorului din baza de date
   - Backend-ul validează codul folosind totp.Validate() care compară codul introdus cu codul calculat server-side
   - Dacă codul este valid, backend-ul creează o sesiune completă pentru utilizator
   - Sesiunea este stocată server-side și utilizatorul primește un cookie de sesiune
   - Frontend-ul actualizează starea aplicației și redirecționează către pagina de profil

### Securitate

#### Măsuri de securitate implementate

Aplicația implementează următoarele măsuri de securitate pentru a proteja datele utilizatorilor și a preveni accesul neautorizat:

1. **Stocarea securizată a secretelor**:
   - Secretele TOTP sunt stocate în baza de date în format text (pentru compatibilitate cu bibliotecile TOTP)
   - Accesul la baza de date este restricționat și protejat prin credențiale sigure
   - Secretele sunt generate folosind generatoare criptografice sigure

2. **Criptarea parolelor**:
   - Parolele sunt hash-uite folosind bcrypt cu cost factor DefaultCost (10)
   - Fiecare parolă are un salt unic generat automat de bcrypt
   - Parolele nu sunt niciodată stocate în format text clar

3. **Protecția sesiunilor**:
   - Sesiunile sunt configurate cu parametri de securitate:
     - HttpOnly: true (prevenirea accesului JavaScript la cookie-uri)
     - SameSite: Lax (protecție împotriva CSRF)
     - MaxAge: 7 zile (expirare automată)
   - Cookie-urile de sesiune sunt semnate pentru a preveni modificarea

4. **Configurare CORS**:
   - Restricționarea originilor permise la frontend-ul aplicației
   - Configurarea headerelor permise și expuse
   - Activarea credențialelor pentru cookie-uri cross-origin

5. **Validarea datelor**:
   - Validarea completă pe server a tuturor datelor primite
   - Sanitizarea inputurilor pentru prevenirea injection attacks
   - Verificarea autorizațiilor pentru fiecare operațiune

6. **Protecție împotriva atacurilor comune**:
   - Rate limiting implicit prin configurația CORS și sesiuni
   - Protecție CSRF prin cookie-uri SameSite
   - Validarea proprietății sesiunii pentru operațiuni sensibile
   - Verificarea parolei pentru dezactivarea MFA

#### Considerații de design pentru securitate

Aplicația a fost proiectată urmând principiul "security by design", cu securitatea integrată în toate nivelurile sistemului:

1. **Defense in Depth** - Multiple straturi de securitate care funcționează independent
2. **Principle of Least Privilege** - Utilizatorii pot modifica doar propriile date
3. **Fail Secure** - În caz de eroare, sistemul refuză accesul
4. **Complete Mediation** - Toate operațiunile sunt verificate și autorizate
5. **Validare pe server** - Nicio decizie de securitate nu se bazează exclusiv pe validări client-side

### Tehnologii Utilizate

#### Backend:
- **Go 1.24** - Limbajul de programare principal
- **Gin Gonic** - Framework web pentru API REST
- **GORM** - ORM pentru operațiuni de bază de date
- **PostgreSQL** - Baza de date relațională
- **bcrypt** - Pentru hash-uirea parolelor
- **pquerna/otp** - Pentru implementarea TOTP
- **Gorilla Sessions** - Pentru gestionarea sesiunilor
- **godotenv** - Pentru gestionarea variabilelor de mediu

#### Frontend:
- **React 18** - Framework JavaScript pentru interfața utilizator
- **TypeScript** - Pentru type safety și dezvoltare mai sigură
- **Vite** - Build tool și development server
- **Tailwind CSS** - Framework CSS pentru styling
- **React Hook Form** - Pentru gestionarea formularelor
- **Zod** - Pentru validarea datelor
- **Zustand** - Pentru state management
- **Axios** - Pentru comunicarea HTTP
- **QRCode** - Pentru generarea codurilor QR
- **React Toastify** - Pentru notificări

### Modalitatea de Rulare

#### Cerințe preliminare

- Go 1.24 sau mai recent
- Node.js 18+ și npm
- PostgreSQL 15+
- Git pentru clonarea repository-urilor

#### Pași pentru rularea cu Docker (Recomandat)

1. **Clonare repository backend:**
   ```bash
   git clone https://github.com/mariusel9911/SSC_backend.git
   cd SSC_backend
   ```

2. **Configurare variabile de mediu:**
   ```bash
   cp .env.example .env
   # Editați .env cu valorile dorite
   ```

3. **Rulare cu Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Clonare și rulare frontend:**
   ```bash
   git clone https://github.com/mariusel9911/SSC_frontend.git
   cd SSC_frontend
   npm install
   npm run dev
   ```

#### Pași pentru rularea în modul dezvoltare

1. **Setup PostgreSQL:**
   ```bash
   # Creați baza de date
   createdb 2fa
   ```

2. **Rulare Backend:**
   ```bash
   cd SSC_backend
   go mod download
   cp .env.example .env
   # Configurați variabilele de mediu în .env
   go run main.go
   ```

3. **Rulare Frontend:**
   ```bash
   cd SSC_frontend
   npm install
   npm run dev
   ```

#### Accesare aplicație:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8000

### Funcionalități Implementate

1. **Autentificare de bază** - Înregistrare și login cu email/parolă
2. **Autentificare cu doi factori** - Configurare și utilizare TOTP
3. **Gestionare profil** - Vizualizare și editare informații utilizator
4. **Activare/Dezactivare 2FA** - Control complet asupra setărilor MFA
5. **Generare QR Code** - Pentru configurarea aplicațiilor de autentificare
6. **Validare în timp real** - Feedback imediat pentru formularele de input
7. **Notificări** - Mesaje de succes și eroare pentru utilizator
8. **Responsive design** - Interfață adaptabilă pentru toate dispozitivele

### Surse de Documentare

În dezvoltarea acestei aplicații, au fost consultate următoarele resurse:

1. **RFC 6238** - TOTP: Time-Based One-Time Password Algorithm
2. **RFC 4226** - HOTP: An HMAC-Based One-Time Password Algorithm
3. **Documentația oficială Go** - https://golang.org/doc/
4. **Biblioteca pquerna/otp** - https://github.com/pquerna/otp
5. **Documentația Gin Gonic** - https://gin-gonic.com/docs/
6. **Documentația React** - https://react.dev/
7. **Ghiduri de securitate OWASP** - Pentru best practices în securitate web
8. **Documentația PostgreSQL** - Pentru optimizarea bazei de date

## Autori

- **Nume Student:** Nistor Marius Ionut
- **Email:** marius.ionut.nistor03@gmail.com
- **Grupa:** 2.2 TI

## Licență

Acest proiect este licențiat sub MIT License
