# Aplicație de Autentificare cu Factori Multipli în Go

## Tema Generală

Această aplicație implementează un sistem de autentificare cu factori multipli (Multi-Factor Authentication - MFA) folosind limbajul de programare Go. Sistemul este compus dintr-un frontend și un backend care interacționează pentru a oferi o soluție completă de autentificare securizată, utilizând multiple metode de verificare pentru a confirma identitatea utilizatorilor.

## Adresele Repository-urilor

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

#### JWT (JSON Web Tokens)

Pentru gestionarea sesiunilor autentificate, aplicația utilizează tehnologia JWT:

- Token-urile JWT sunt semnate digital pentru a asigura integritatea informațiilor
- Conțin informații despre utilizator (claims) într-un format standardizat
- Permit implementarea unui sistem de autentificare stateless
- Includ informații despre perioada de valabilitate a sesiunii
- Pot fi reînnoite folosind un sistem de refresh tokens

### Arhitectură și Design

#### Arhitectura Generală

Aplicația este structurată folosind un model client-server, cu separare clară între frontend și backend:

1. **Frontend** - interfața cu utilizatorul, dezvoltată cu tehnologii web moderne și găzduită în repository-ul SSC_frontend
2. **Backend** - serviciul API care gestionează logica de autentificare, dezvoltat în Go și găzduit în repository-ul SSC_backend

Comunicarea între cele două componente se realizează prin API RESTful.

#### Principii arhitecturale

Aplicația implementează următoarele principii de arhitectură:

1. **Separation of Concerns (SoC)** - Separarea clară a responsabilităților între componente, cu frontend-ul responsabil pentru interfața utilizator și backend-ul pentru logica de business și securitate.

2. **Backend for Frontend (BFF)** - Un model arhitectural în care backend-ul este special conceput pentru a servi nevoile frontend-ului, optimizând comunicarea între cele două componente.

3. **Arhitectură pe straturi** - Backend-ul este structurat în straturi logice (prezentare, servicii, date) pentru a separa responsabilitățile și a facilita testarea și mentenanța.

4. **Model API-First** - Dezvoltarea a început cu definirea clară a API-ului, care a servit ca un contract între frontend și backend.

5. **Stateless Authentication** - Utilizarea token-urilor JWT pentru a implementa o autentificare fără stare, îmbunătățind scalabilitatea sistemului.

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
   - `AuthMiddleware`: Verifică validitatea token-urilor JWT pentru rutele protejate
   - `RateLimiter`: Limitează numărul de cereri pentru a preveni atacurile de forță brută

4. **Generator TOTP** - creează și validează codurile temporare folosind algoritmi standardizați:
   - Implementat folosind biblioteca `pquerna/otp` care respectă standardul RFC 6238
   - Generează secrete unice pentru fiecare utilizator
   - Validează codurile cu toleranță temporală pentru a compensa decalajele de timp între dispozitive

5. **Servicii** - Implementează logica de business:
   - `UserService`: Gestionează operațiunile legate de utilizatori
   - `TokenService`: Creează și validează token-urile JWT
   - `OTPService`: Gestionează operațiunile legate de autentificarea cu doi factori

6. **Stocare date** - gestionează persistența informațiilor:
   - Utilizează o bază de date PostgreSQL pentru stocarea datelor utilizatorilor
   - Implementează Repository Pattern pentru abstractizarea operațiunilor pe baza de date
   - Utilizează criptare pentru stocarea securizată a secretelor TOTP

7. **Utilitare**:
   - Logare structurată pentru monitorizarea activității
   - Gestionare configurație din variabile de mediu
   - Module helper pentru validarea datelor

Biblioteca principală utilizată pentru implementarea TOTP este bazată pe standardul RFC 6238, care definește algoritmi OATH-TOTP pentru generarea codurilor temporare.

#### Structura directorului backend

```
SSC_backend/
├── cmd/
│   └── server/
│       └── main.go                  # Punctul de intrare al aplicației
├── internal/
│   ├── api/
│   │   ├── controllers/             # Controllere pentru gestionarea cererilor HTTP
│   │   │   ├── auth_controller.go
│   │   │   └── otp_controller.go
│   │   ├── middleware/              # Middleware pentru procesarea cererilor
│   │   │   ├── auth_middleware.go
│   │   │   └── rate_limiter.go
│   │   └── routes/
│   │       └── routes.go            # Definirea rutelor API
│   ├── config/
│   │   └── config.go                # Configurația aplicației
│   ├── models/
│   │   ├── user.go                  # Modelul de date pentru utilizatori
│   │   └── token.go                 # Modelul de date pentru token-uri
│   ├── repository/
│   │   ├── user_repository.go       # Operațiuni pe baza de date pentru utilizatori
│   │   └── otp_repository.go        # Operațiuni pentru gestionarea secretelor TOTP
│   └── services/
│       ├── auth_service.go          # Serviciu pentru autentificare
│       ├── otp_service.go           # Serviciu pentru autentificare cu doi factori
│       └── token_service.go         # Serviciu pentru gestionarea token-urilor JWT
├── pkg/
│   ├── crypto/                      # Utilitar pentru criptare/decriptare
│   ├── logger/                      # Utilitar pentru logare
│   └── validator/                   # Utilitar pentru validarea datelor
├── go.mod
└── go.sum
```

#### Componente Frontend

Frontend-ul este dezvoltat utilizând tehnologii moderne pentru a oferi o experiență de utilizare fluidă și intuitivă. Principalele componente includ:

1. **Framework UI** - Implementat folosind tehnologii moderne de frontend:
   - Framework React pentru construirea interfețelor reutilizabile
   - State management pentru gestionarea stării aplicației
   - Rutare client-side pentru navigarea fără refresh

2. **Module de autentificare**:
   - `AuthProvider`: Component React care gestionează starea de autentificare
   - `ProtectedRoute`: Component pentru protejarea rutelor care necesită autentificare
   - `TokenManager`: Serviciu pentru gestionarea token-urilor JWT în localStorage/sessionStorage

3. **Componente de interfață**:
   - `RegisterForm`: Formular pentru înregistrarea utilizatorilor noi
   - `LoginForm`: Formular pentru autentificarea cu utilizator și parolă
   - `OTPSetup`: Interfață pentru configurarea autentificării cu doi factori
   - `OTPValidator`: Component pentru introducerea și validarea codurilor TOTP
   - `UserProfile`: Interfață pentru gestionarea profilului utilizatorului

4. **Servicii**:
   - `ApiService`: Gestionează comunicarea cu backend-ul prin API RESTful
   - `QrCodeGenerator`: Generează coduri QR pentru scanarea cu aplicații de autentificare
   - `NotificationService`: Afișează notificări și mesaje pentru utilizator

5. **Utilitare**:
   - Module pentru validarea formularelor
   - Interceptori pentru cereri HTTP
   - Module pentru internaționalizare (i18n)

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
│   ├── index.html
│   └── assets/
├── src/
│   ├── components/
│   │   ├── auth/
│   │   │   ├── LoginForm.js
│   │   │   ├── RegisterForm.js
│   │   │   ├── OTPSetup.js
│   │   │   └── OTPValidator.js
│   │   ├── common/
│   │   │   ├── Button.js
│   │   │   ├── Input.js
│   │   │   └── Notification.js
│   │   └── profile/
│   │       └── UserProfile.js
│   ├── contexts/
│   │   └── AuthContext.js
│   ├── hooks/
│   │   ├── useAuth.js
│   │   └── useApi.js
│   ├── pages/
│   │   ├── HomePage.js
│   │   ├── LoginPage.js
│   │   ├── RegisterPage.js
│   │   ├── ProfilePage.js
│   │   └── OTPPage.js
│   ├── services/
│   │   ├── api.service.js
│   │   ├── auth.service.js
│   │   └── qrcode.service.js
│   ├── utils/
│   │   ├── validators.js
│   │   └── tokenManager.js
│   ├── App.js
│   └── index.js
├── package.json
└── package-lock.json
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
       │                   │   Verifică JWT    │                   │
       │                   │───────────────────>                   │
       │                   │                   │                   │
       │                   │  JWT invalid/lipsă│                   │
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
       │                   │  Generează JWT    │                   │
       │                   │<───────────────────                   │
       │                   │                   │                   │
       │  Acces acordat    │                   │                   │
       │<──────────────────│                   │                   │
       │                   │                   │                   │
```

#### Procese detaliate

1. **Înregistrare**:
   - Utilizatorul accesează pagina de înregistrare și completează formularul cu datele personale (email, parolă, etc.)
   - Frontend-ul validează datele introduse (formate corecte, parolă suficient de puternică)
   - Datele validate sunt trimise către backend prin API
   - Backend-ul verifică unicitatea email-ului în baza de date
   - Parola este hash-uită folosind funcții criptografice sigure (bcrypt sau Argon2)
   - Datele utilizatorului sunt stocate în baza de date
   - Backend-ul generează un JWT pentru sesiunea inițială
   - Utilizatorul este redirecționat către pagina de profil sau dashboard

2. **Activare MFA**:
   - Utilizatorul autentificat accesează pagina de profil și selectează opțiunea de activare MFA
   - Frontend-ul trimite o cerere către backend pentru inițierea procesului
   - Backend-ul generează un secret unic pentru utilizator
   - "Secret"-ul este convertit într-un format URI compatibil cu aplicațiile de autentificare
   - Backend-ul generează datele necesare pentru crearea codului QR
   - Frontend-ul primește datele și afișează codul QR pentru utilizator
   - Utilizatorul scanează codul QR cu o aplicație de autentificare (Google Authenticator, Authy, etc.)
   - Aplicația de autentificare începe să genereze coduri TOTP pe baza "secret"-ului
   - Utilizatorul introduce codul curent afișat în aplicație în interfața web
   - Frontend-ul trimite codul către backend pentru validare
   - Backend-ul verifică codul folosind un "secret" stocat
   - Dacă validarea reușește, MFA este activat în profilul utilizatorului
   - "Secret"-ul este stocat criptat în baza de date pentru utilizări viitoare

3. **Autentificare cu MFA**:
   - Utilizatorul accesează pagina de autentificare și introduce email și parolă (primul factor)
   - Frontend-ul trimite credențialele către backend
   - Backend-ul validează credențialele și verifică dacă utilizatorul are MFA activat
   - Pentru utilizatorii cu MFA activat, backend-ul răspunde cu un token temporar și solicită al doilea factor
   - Frontend-ul afișează ecranul pentru introducerea codului TOTP
   - Utilizatorul consultă aplicația de autentificare și introduce codul generat
   - Frontend-ul trimite codul împreună cu token-ul temporar către backend
   - Backend-ul recuperează "secret"-ul utilizatorului din baza de date și validează codul TOTP
   - Backend-ul verifică validitatea codului, luând în considerare și mici decalaje temporale
   - Dacă codul este valid, backend-ul generează un JWT complet pentru sesiune
   - Frontend-ul stochează token-ul JWT și redirectează utilizatorul către dashboard
   - Toate cererile ulterioare către API vor include token-ul JWT pentru autorizare

### Securitate

#### Măsuri de securitate implementate

Aplicația implementează următoarele măsuri de securitate pentru a proteja datele utilizatorilor și a preveni accesul neautorizat:

1. **Stocarea securizată a secretelor**:
   - Secretele TOTP sunt criptate folosind AES-256 înainte de stocare în baza de date
   - Cheile de criptare sunt rotite periodic și stocate separat de datele criptate
   - Este implementat un mecanism de auditing pentru monitorizarea accesului la secrete

2. **Criptarea parolelor**:
   - Parolele sunt hash-uite folosind algoritmi moderni (Argon2 sau bcrypt)
   - Sunt utilizați factori de cost adecvați pentru a face atacurile de forță brută impracticabile
   - Salt-uri unice sunt generate pentru fiecare parolă pentru a preveni atacurile cu tabele rainbow

3. **Rate limiting**:
   - Limitarea numărului de încercări de autentificare pentru prevenirea atacurilor de forță brută
   - Implementarea de întârzieri progresive pentru încercări eșuate multiple
   - Monitorizarea și blocarea temporară a adreselor IP suspecte

4. **Sesiuni securizate**:
   - Utilizarea de token-uri JWT cu semnătură digitală pentru autentificare
   - Durata limitată de viață a token-urilor pentru a reduce fereastra de oportunitate în caz de compromitere
   - Implementarea unui sistem de reîmprospătare a token-urilor pentru experiență optimă de utilizare
   - Mecanisme de invalidare a token-urilor pentru delogare sau în caz de compromitere

5. **Protecție CSRF**:
   - Implementarea de token-uri anti-CSRF pentru formulare
   - Verificarea headerelor Origin și Referer pentru cereri sensibile
   - Utilizarea de cookie-uri cu atributele SameSite și HttpOnly

6. **Alte măsuri**:
   - Headerele HTTP de securitate (Content-Security-Policy, X-XSS-Protection, etc.)
   - Sanitizarea input-urilor pentru prevenirea atacurilor de tip injecție
   - Validarea strictă a datelor pe server, indiferent de validarea din client
   - Utilizarea HTTPS pentru toate comunicațiile
   - Logging și auditing pentru detectarea tentativelor de intruziune

#### Considerații de design pentru securitate

Aplicația a fost proiectată urmând principiul "security by design", cu securitatea integrată în toate nivelurile sistemului:

1. **Defense in Depth** - Multiple straturi de securitate care funcționează independent
2. **Principle of Least Privilege** - Acordarea doar a permisiunilor minime necesare pentru funcționare
3. **Fail Secure** - În caz de eroare, sistemul defaultează la starea sigură (acces refuzat)
4. **Complete Mediation** - Toate accesele la resurse sunt verificate complet
5. **Validare pe server** - Nicio decizie de securitate nu se bazează exclusiv pe validări în client

### Modalitatea de Rulare

#### Cerințe preliminare

- Go (versiunea 1.16 sau mai recentă)
- Node.js și npm (pentru frontend)
- Bază de date PostgreSQL (sau altă bază configurată în fișierul de configurare)

#### Pași pentru rularea Backend

1. Clonare repository:
   ```
   git clone https://github.com/mariusel9911/SSC_backend.git
   cd SSC_backend
   ```

2. Instalare dependențe:
   ```
   go mod download
   ```

3. Configurare variabile de mediu:
   Creați un fișier `.env` în directorul principal și configurați variabilele necesare (consultați fișierul `.env.example` pentru un exemplu)

4. Rulare aplicație:
   ```
   go run main.go
   ```

Backend-ul va rula implicit pe portul 8080 sau pe portul specificat în variabilele de mediu.

#### Pași pentru rularea Frontend

1. Clonare repository:
   ```
   git clone https://github.com/mariusel9911/SSC_frontend.git
   cd SSC_frontend
   ```

2. Instalare dependențe:
   ```
   npm install
   ```

3. Configurare variabile de mediu:
   Creați un fișier `.env` în directorul principal cu variabilele necesare (consultați fișierul `.env.example` pentru un exemplu)

4. Rulare aplicație:
   ```
   npm start
   ```

Frontend-ul va rula implicit pe portul 3000.

### Surse de Documentare

În dezvoltarea acestei aplicații, au fost consultate următoarele resurse:

1. RFC 6238 - TOTP: Time-Based One-Time Password Algorithm
2. RFC 4226 - HOTP: An HMAC-Based One-Time Password Algorithm
3. Documentația oficială Go: https://golang.org/doc/
4. Biblioteca Go pentru TOTP: https://github.com/pquerna/otp
5. Ghiduri de securitate OWASP pentru implementarea MFA
6. Documentație pentru best practices în implementarea JWT pentru sesiuni

## Autori

- Nume Student: Nistor Marius Ionut
- Email: marius.ionut.nistor03@gmail.com
- Grupa: 2.2 TI

## Licență

Acest proiect este licențiat sub MIT LICENCE
