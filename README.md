# Objective

Synchronization of workouts from endomondo to strava.

# Current status

Creating backup of endomondo workouts, importing to strava is not done yet.

# Why?

Endomondo will die till the end of 2020 (Press F to pay respect)

# Why not tapirrik

Tapiriik doesn't work

# Booting up

1. Rename/copy .env.sample to .env. ```cp .env.sample .env```
2. Fill with your endomondo data email/pass
3. Set start backup and end backup date
4. Run ```make all```