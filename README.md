# Endomondo2Strava

[![Build Status](https://travis-ci.com/michalq/endomondo2strava.svg?branch=master)](https://travis-ci.com/michalq/endomondo2strava)

### Objective

Synchronization of workouts from endomondo to strava.

### Why?

Endomondo will die till the end of 2020 (Press F to pay respect)

### Why not tapiriik?

Tapiriik doesn't work, google says that not just for me…

# Configuration

Configuration is located in ```.env``` file, which first have to be created or copy from ```.env.sample```.

Possible values:

- **ENDOMONDO_EMAIL** Endomondo email
    
    ```ENDOMONDO_EMAIL=your@email```
- **ENDOMONDO_PASS** Endomondo password
    
    ```ENDOMONDO_PASS=xyz```
- **STRAVA_CLIENT_ID** Strava client id
    
    ```STRAVA_CLIENT_SECRET=ABCD```
- **STRAVA_CLIENT_SECRET** Strava client secret
    
    ```STRAVA_CLIENT_ID=abcd```
- **ENDOMONDO_EXPORT_FORMAT** Format in which workouts will be downloaded from endomondo. 
    
    Possible values <GPX, TCX>.
    
    ```ENDOMONDO_EXPORT_FORMAT=GPX```
- **STEP** You can skip downloading part if you already have downloaded workouts by passing here just upload.
    
    Possible values <export, import>
    - export: runs only export from endomondo
    - import: runs only import to strava
    ```
    // Does full synchro
    STEP=export,import
    // Only export
    STEP=export
    // Only import
    STEP=import
    ```

# Booting up

### Preconditions

1. GoLang installed, version 1.15

### Steps to run

1. Rename/copy .env.sample to .env. ```cp .env.sample .env```
2. Fill with your endomondo data email/pass
3. Generate your strava client-id/client-secret [here](https://www.strava.com/settings/api) and copy clientId and clientSecret to configuration,
4. Set start backup and end backup date,
5. Select type of export gpx or tcx,
5. Run ```make all```*

*With the first run you will have to authorize yourself in strava by opening link, clicking Authorize and then copying code from url. Authorization data are saved in db so next run doesn't require that step.

Example result output:
```
----------------------------
Found workouts          574
----------------------------
Downloaded details      574
----------------------------
Found photos            3
----------------------------
Downloaded workouts     574
----------------------------
Imported workouts       100
----------------------------
```

### Results

All exported workouts you will wind under ```./tmp/workouts```, in ```./tmp/db.sqlite``` you can find database with export/import status.

### Limitations

Strava allows only for 100 requests to api per 15 minutes. If you have more than 100 workouts you will have to run script few times waiting 15 minutes after each run. In database is saved information which workouts were uploaded, so each run will send another not imported workouts. If you want to run script many times make sure to set ```STEP=import``` to skip exporting each time from endomondo.

Endomondo doesn't have limitations so probably you will be able to download all workouts on first session. However endomondo firewall can block you if you made to much requests, in that case it is good to run export few times with different start and end date.

# Dependencies

- Strava API/v3
    
    For uploading to strava is used official strava API/v3, which is documented here: https://developers.strava.com/docs/reference/.
- Endomondo 
    
    Endomondo doesn't have official API thus it's used session API from www, email/pass is required.

# Q

1. Does it duplicate my workouts?

    It tracks in sqlite what was imported, and what not, so shouldn't. Strava seems to have also a security check for uploading same workouts many times. 

# Known problems

1. Simple workouts without tracking points are treated as "empty file by strava"

# TODO

1. Dockerize it

# Changelog

## 0.2.0
 - Saving in DB title, description and hashtags from endomondo
 - Description is saved in format: "{endomondoDescription} {hashtags} (Endomondo id {endomondiId})"
 - Photos are saved in local database as a link (have to be queried manually),
 - Added report at the end of eport/import