/**
 * Appvia Kore API
 * Kore API provides the frontend API for the Appvia Kore (kore.appvia.io)
 *
 * The version of the OpenAPI document: 0.0.1
 * Contact: info@appvia.io
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';
import V1Ownership from './V1Ownership';

/**
 * The V1SecurityResourceOverview model module.
 * @module model/V1SecurityResourceOverview
 * @version 0.0.1
 */
class V1SecurityResourceOverview {
    /**
     * Constructs a new <code>V1SecurityResourceOverview</code>.
     * @alias module:model/V1SecurityResourceOverview
     */
    constructor() { 
        
        V1SecurityResourceOverview.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1SecurityResourceOverview</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1SecurityResourceOverview} obj Optional instance to populate.
     * @return {module:model/V1SecurityResourceOverview} The populated <code>V1SecurityResourceOverview</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1SecurityResourceOverview();

            if (data.hasOwnProperty('lastChecked')) {
                obj['lastChecked'] = ApiClient.convertToType(data['lastChecked'], 'String');
            }
            if (data.hasOwnProperty('openIssueCounts')) {
                obj['openIssueCounts'] = ApiClient.convertToType(data['openIssueCounts'], {'String': 'Number'});
            }
            if (data.hasOwnProperty('overallStatus')) {
                obj['overallStatus'] = ApiClient.convertToType(data['overallStatus'], 'String');
            }
            if (data.hasOwnProperty('resource')) {
                obj['resource'] = V1Ownership.constructFromObject(data['resource']);
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getLastChecked() {
        return this.lastChecked;
    }

    /**
     * @param {String} lastChecked
     */
    setLastChecked(lastChecked) {
        this['lastChecked'] = lastChecked;
    }
/**
     * @return {Object.<String, Number>}
     */
    getOpenIssueCounts() {
        return this.openIssueCounts;
    }

    /**
     * @param {Object.<String, Number>} openIssueCounts
     */
    setOpenIssueCounts(openIssueCounts) {
        this['openIssueCounts'] = openIssueCounts;
    }
/**
     * @return {String}
     */
    getOverallStatus() {
        return this.overallStatus;
    }

    /**
     * @param {String} overallStatus
     */
    setOverallStatus(overallStatus) {
        this['overallStatus'] = overallStatus;
    }
/**
     * @return {module:model/V1Ownership}
     */
    getResource() {
        return this.resource;
    }

    /**
     * @param {module:model/V1Ownership} resource
     */
    setResource(resource) {
        this['resource'] = resource;
    }

}

/**
 * @member {String} lastChecked
 */
V1SecurityResourceOverview.prototype['lastChecked'] = undefined;

/**
 * @member {Object.<String, Number>} openIssueCounts
 */
V1SecurityResourceOverview.prototype['openIssueCounts'] = undefined;

/**
 * @member {String} overallStatus
 */
V1SecurityResourceOverview.prototype['overallStatus'] = undefined;

/**
 * @member {module:model/V1Ownership} resource
 */
V1SecurityResourceOverview.prototype['resource'] = undefined;






export default V1SecurityResourceOverview;
