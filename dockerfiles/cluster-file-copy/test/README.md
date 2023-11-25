## push

```sh
CLUSTER=main START_MESSAGE=~fits-fz#Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M16_0005.fits.fz~qiu scalebox app create

CLUSTER=main START_MESSAGE=~fits-fz#Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M16_0005.fits.fz~p419 scalebox app create
```
## pull 
```sh
CLUSTER=main START_MESSAGE=p419~fits-2#Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M16_0005.fits.fz~ scalebox app create

CLUSTER=qiu START_MESSAGE=main~fits-fz#Dec+6007_09_03/20221019/Dec+6007_09_03_arcdrift-M16_0005.fits.fz~ scalebox app create
```
