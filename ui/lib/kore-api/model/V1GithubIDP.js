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

/**
 * The V1GithubIDP model module.
 * @module model/V1GithubIDP
 * @version 0.0.1
 */
class V1GithubIDP {
    /**
     * Constructs a new <code>V1GithubIDP</code>.
     * @alias module:model/V1GithubIDP
     * @param clientID {String} 
     * @param clientSecret {String} 
     * @param orgs {Array.<String>} 
     */
    constructor(clientID, clientSecret, orgs) { 
        
        V1GithubIDP.initialize(this, clientID, clientSecret, orgs);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, clientID, clientSecret, orgs) { 
        obj['clientID'] = clientID;
        obj['clientSecret'] = clientSecret;
        obj['orgs'] = orgs;
    }

    /**
     * Constructs a <code>V1GithubIDP</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1GithubIDP} obj Optional instance to populate.
     * @return {module:model/V1GithubIDP} The populated <code>V1GithubIDP</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1GithubIDP();

            if (data.hasOwnProperty('clientID')) {
                obj['clientID'] = ApiClient.convertToType(data['clientID'], 'String');
            }
            if (data.hasOwnProperty('clientSecret')) {
                obj['clientSecret'] = ApiClient.convertToType(data['clientSecret'], 'String');
            }
            if (data.hasOwnProperty('orgs')) {
                obj['orgs'] = ApiClient.convertToType(data['orgs'], ['String']);
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getClientID() {
        return this.clientID;
    }

    /**
     * @param {String} clientID
     */
    setClientID(clientID) {
        this['clientID'] = clientID;
    }
/**
     * @return {String}
     */
    getClientSecret() {
        return this.clientSecret;
    }

    /**
     * @param {String} clientSecret
     */
    setClientSecret(clientSecret) {
        this['clientSecret'] = clientSecret;
    }
/**
     * @return {Array.<String>}
     */
    getOrgs() {
        return this.orgs;
    }

    /**
     * @param {Array.<String>} orgs
     */
    setOrgs(orgs) {
        this['orgs'] = orgs;
    }

}

/**
 * @member {String} clientID
 */
V1GithubIDP.prototype['clientID'] = undefined;

/**
 * @member {String} clientSecret
 */
V1GithubIDP.prototype['clientSecret'] = undefined;

/**
 * @member {Array.<String>} orgs
 */
V1GithubIDP.prototype['orgs'] = undefined;






export default V1GithubIDP;

