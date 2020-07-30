import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'timestamp'
})
export class TimestampPipe implements PipeTransform {
    transform(item: string, showTimestamp: boolean): string {
        const regex: RegExp = new RegExp('^[^ ]+');
        if (!showTimestamp) {
            return item.replace(regex, '').trim();
        } else {
            return item.replace(regex, (match) => `<span class="timestamp">${match}</span>`);
        }
    }
}
