/**
 * Kore API
 * Kore API provides the frontend API (kore.appvia.io)
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
 * The V1AllocationSpec model module.
 * @module model/V1AllocationSpec
 * @version 0.0.1
 */
class V1AllocationSpec {
    /**
     * Constructs a new <code>V1AllocationSpec</code>.
     * @alias module:model/V1AllocationSpec
     * @param name {String} 
     * @param resource {module:model/V1Ownership} 
     * @param summary {String} 
     * @param teams {Array.<String>} 
     */
    constructor(name, resource, summary, teams) { 
        
        V1AllocationSpec.initialize(this, name, resource, summary, teams);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, name, resource, summary, teams) { 
        obj['name'] = name;
        obj['resource'] = resource;
        obj['summary'] = summary;
        obj['teams'] = teams;
    }

    /**
     * Constructs a <code>V1AllocationSpec</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1AllocationSpec} obj Optional instance to populate.
     * @return {module:model/V1AllocationSpec} The populated <code>V1AllocationSpec</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1AllocationSpec();

            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('resource')) {
                obj['resource'] = V1Ownership.constructFromObject(data['resource']);
            }
            if (data.hasOwnProperty('summary')) {
                obj['summary'] = ApiClient.convertToType(data['summary'], 'String');
            }
            if (data.hasOwnProperty('teams')) {
                obj['teams'] = ApiClient.convertToType(data['teams'], ['String']);
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getName() {
        return this.name;
    }

    /**
     * @param {String} name
     */
    setName(name) {
        this['name'] = name;
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
/**
     * @return {String}
     */
    getSummary() {
        return this.summary;
    }

    /**
     * @param {String} summary
     */
    setSummary(summary) {
        this['summary'] = summary;
    }
/**
     * @return {Array.<String>}
     */
    getTeams() {
        return this.teams;
    }

    /**
     * @param {Array.<String>} teams
     */
    setTeams(teams) {
        this['teams'] = teams;
    }

}

/**
 * @member {String} name
 */
V1AllocationSpec.prototype['name'] = undefined;

/**
 * @member {module:model/V1Ownership} resource
 */
V1AllocationSpec.prototype['resource'] = undefined;

/**
 * @member {String} summary
 */
V1AllocationSpec.prototype['summary'] = undefined;

/**
 * @member {Array.<String>} teams
 */
V1AllocationSpec.prototype['teams'] = undefined;






export default V1AllocationSpec;

