import { DOCUMENT } from '@angular/common';
import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { Inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class HttpInterceptorService implements HttpInterceptor {

    constructor(
        @Inject(DOCUMENT) private document: Document
    ) { }

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        const baseHref = this.getBaseHref();

        // If the base href value is not the default "/" and the requested resource is not an absolute URL,
        // we have to adjust the requested resource URL, in order to account for the base href itself
        if (baseHref !== '/' && !req.url.match(/http(?:s){0,1}:\/\//)) {
            let url = req.url;
            if (req.url.startsWith('/')) {
                // If the requested URL starts with a slash, we will remove it
                url = req.url.slice(1);
            }

            req = req.clone({
                url: `${baseHref}${url}`
            });
        }

        return next.handle(req);
    }

    /**
     * An helper method to get base href. If, for whatever reason, we cannot find a
     * "base" tag with an "href" attribute, this method will return "/".
     * Result of this function is a string that will always end with a slash.
     * @returns the base href
     */
    private getBaseHref() {
        const baseHrefTag = this.document.querySelector('base[href]');
        if (baseHrefTag) {
            let href = baseHrefTag.getAttribute('href') || '/';
            if (!href.endsWith('/')) {
                href += '/';
            }

            return href;
        }

        return '/';
    }

}
