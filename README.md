# About

### Objective

Synchronization of workouts from endomondo to strava.

### Current status

Creating backup of endomondo workouts, importing to strava is not done yet.

### Why?

Endomondo will die till the end of 2020 (Press F to pay respect)

### Why not tapiriik?

Tapiriik doesn't work

# Booting up

1. Rename/copy .env.sample to .env. ```cp .env.sample .env```
2. Fill with your endomondo data email/pass
3. Set start backup and end backup date
4. Run ```make all```

# Environments

### Endomondo login data

```
ENDOMONDO_EMAIL=your@email
ENDOMONDO_PASS=xyz
```

### Strava OAuth2.0 client id and client secret

```
STRAVA_CLIENT_ID=abcd
STRAVA_CLIENT_SECRET=ABCD
```

### Export/import

Starting point to search workouts to export
```
START_AT=2020-01-01
```

Ending point to search workouts to export
```
END_AT=2020-11-01
```

Format in which workouts will be downloaded from endomondo. 
Possible values <GPX, TCX>.
```
ENDOMONDO_EXPORT_FORMAT=GPX
```

You can skip downloading part if you already have downloaded workouts by passing here just upload.
Possible values <export, import>
- export: runs only export from endomondo
- import: runs only import to strava

```
// Does full synchro
STEP=export,import
// Only export
STEP=export
```

# Used API's

### Strava

For uploading to strava is used official strava API/v3, which is documented here: https://developers.strava.com/docs/reference/.

### Endomondo

Endomondo doesn't have official API thus it's used session API from www, thus email/pass is required.