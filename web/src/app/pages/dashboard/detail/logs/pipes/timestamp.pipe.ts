import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';
import 'moment-timezone';

@Pipe({
    name: 'timestamp'
})
export class TimestampPipe implements PipeTransform {
    transform(item: string, showTimestamp: boolean): string {
        const timestampRegex: RegExp = new RegExp('^[^ ]+');
        if (!showTimestamp) {
            return item.replace(timestampRegex, '').trim();
        } else {
            return item.replace(timestampRegex, (match) => {
                const formattedTimestamp: string = moment(match).format('MMMM Do YYYY, h:mm:ss a z');
                const zone: string = moment.tz.guess();
                const timezone: string = moment.tz(zone).format('z');
                return `<span style="color:#777777;font-weight:500;">${formattedTimestamp} ${timezone}</span>`;
            });
        }
    }
}
