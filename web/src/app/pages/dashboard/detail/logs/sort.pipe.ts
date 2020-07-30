import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';

@Pipe({
    name: 'sort'
})
export class SortPipe implements PipeTransform {
    transform(logs: Log[], order: string): Log[] {
        console.log(order);
        if (!order) {
            return logs;
        }
        const comparator = (a: Log, b: Log) => {
            return Date.parse(a.timestamp) - Date.parse(b.timestamp);
        }
        if (order === 'asc') {
            return logs.sort(comparator);
        } else {
            return logs.sort(comparator).reverse();
        }
    }
}
