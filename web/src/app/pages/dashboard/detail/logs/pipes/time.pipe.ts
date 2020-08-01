import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';
import * as moment from 'moment';

// number of nanoseconds in a millisecond
const nanoMilli = 1000000;

@Pipe({
    name: 'time'
})
export class TimePipe implements PipeTransform {
    transform(items: Log[], dateFrom: Date, timeFrom: string, dateTo: Date, timeTo: string): Log[] {
        if (!dateFrom) {
            if (!timeFrom) {
                return items;
            }
        }
        if (!timeFrom) {
            timeFrom = '00:00';
        }
        const from = moment(`${dateFrom.getFullYear()}-${dateFrom.getMonth() + 1}-${dateFrom.getDate()} ${timeFrom}`);
        let to = moment();
        if (!dateTo) {
            if (timeTo) {
                const date = new Date();
                to = moment(`${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${timeTo}`);
            }
        } else {
            if (!timeTo) {
                to = moment(`${dateTo.getFullYear()}-${dateTo.getMonth() + 1}-${dateTo.getDate()}`);
            } else {
                to = moment(`${dateTo.getFullYear()}-${dateTo.getMonth() + 1}-${dateTo.getDate()} ${timeTo}`);
            }
        }
        return items.filter(item => {
            return ((from.utc().valueOf() * nanoMilli) <= item.timestamp) && ((to.utc().valueOf() * nanoMilli) >= item.timestamp);
        });
    }
}
