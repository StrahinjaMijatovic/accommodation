# accommodation

○ auth-service - Čuva kredencijale korisnika i njihove uloge u sistemu. Zadužen za
proces registracije i prijave korisnika.
○ profile-service - Sadrži informacije o osnovnim podacima korisnika kao što su
ime, prezime, pol, starost, imejl itd.
○ accommodations-service - Upravlja osnovnim informacijama o smeštaju (naziv,
opis, slike itd)
○ reservations-service - Kontroliše periode dostupnosti i cene smeštaja, kao i sve
rezervacije
○ recommendations-service - Podržava operacije preporuke smeštaja za goste
○ notifications-service - Upravlja skladištenjem i slanjem notifikacija korisnicima

Baze podataka:
    - auth-service - mongoDB
    - profile-service - mongoDB
    - accommodations-service - cassandra
    - reservations-service - neo4j
    - recommendations-service
    - notifications-service