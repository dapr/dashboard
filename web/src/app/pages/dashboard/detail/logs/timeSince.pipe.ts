import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';

const nano = 1000000000;
const nanoMilli = 1000000;

@Pipe({
    name: 'timeSince'
})
export class TimeSincePipe implements PipeTransform {
    transform(items: Log[], since: number, sinceUnit: string): Log[] {
        if (!items) { return []; }
        if (!since || !sinceUnit) { return items; }
        return items.filter(item => {
            return ((Date.now() * nanoMilli) - item.timestamp) < getNanoseconds(since, sinceUnit);
        });
    }
}

// Get number of nanoseconds since given time
function getNanoseconds(since: number, sinceUnit: string): number {
    // number of nanoseconds in a second
    switch (sinceUnit) {
        case 'seconds':
            return since * nano;
        case 'minutes':
            return since * 60 * nano;
        case 'hours':
            return since * 60 * 60 * nano;
        case 'days':
            return since * 60 * 60 * 24 * nano;
        default:
            throw new Error('Unsupported time since unit');
    }
}
