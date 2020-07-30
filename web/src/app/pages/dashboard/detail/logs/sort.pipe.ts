import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';

@Pipe({
    name: 'sort'
})
export class SortPipe implements PipeTransform {
    transform(logs: Log[], order: string): Log[] {
        if (!order) {
            return logs;
        }
        const comparator = (a: Log, b: Log) => {
            return a.timestamp - b.timestamp;
        };
        if (order === 'asc') {
            return logs.sort(comparator);
        } else {
            return logs.sort(comparator).reverse();
        }
    }
}
