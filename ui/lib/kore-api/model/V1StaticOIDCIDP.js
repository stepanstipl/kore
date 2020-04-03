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
 * The V1StaticOIDCIDP model module.
 * @module model/V1StaticOIDCIDP
 * @version 0.0.1
 */
class V1StaticOIDCIDP {
    /**
     * Constructs a new <code>V1StaticOIDCIDP</code>.
     * @alias module:model/V1StaticOIDCIDP
     * @param clientID {String} 
     * @param clientScopes {Array.<String>} 
     * @param clientSecret {String} 
     * @param issuer {String} 
     * @param userClaims {Array.<String>} 
     */
    constructor(clientID, clientScopes, clientSecret, issuer, userClaims) { 
        
        V1StaticOIDCIDP.initialize(this, clientID, clientScopes, clientSecret, issuer, userClaims);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, clientID, clientScopes, clientSecret, issuer, userClaims) { 
        obj['clientID'] = clientID;
        obj['clientScopes'] = clientScopes;
        obj['clientSecret'] = clientSecret;
        obj['issuer'] = issuer;
        obj['userClaims'] = userClaims;
    }

    /**
     * Constructs a <code>V1StaticOIDCIDP</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1StaticOIDCIDP} obj Optional instance to populate.
     * @return {module:model/V1StaticOIDCIDP} The populated <code>V1StaticOIDCIDP</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1StaticOIDCIDP();

            if (data.hasOwnProperty('clientID')) {
                obj['clientID'] = ApiClient.convertToType(data['clientID'], 'String');
            }
            if (data.hasOwnProperty('clientScopes')) {
                obj['clientScopes'] = ApiClient.convertToType(data['clientScopes'], ['String']);
            }
            if (data.hasOwnProperty('clientSecret')) {
                obj['clientSecret'] = ApiClient.convertToType(data['clientSecret'], 'String');
            }
            if (data.hasOwnProperty('issuer')) {
                obj['issuer'] = ApiClient.convertToType(data['issuer'], 'String');
            }
            if (data.hasOwnProperty('userClaims')) {
                obj['userClaims'] = ApiClient.convertToType(data['userClaims'], ['String']);
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
     * @return {Array.<String>}
     */
    getClientScopes() {
        return this.clientScopes;
    }

    /**
     * @param {Array.<String>} clientScopes
     */
    setClientScopes(clientScopes) {
        this['clientScopes'] = clientScopes;
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
     * @return {String}
     */
    getIssuer() {
        return this.issuer;
    }

    /**
     * @param {String} issuer
     */
    setIssuer(issuer) {
        this['issuer'] = issuer;
    }
/**
     * @return {Array.<String>}
     */
    getUserClaims() {
        return this.userClaims;
    }

    /**
     * @param {Array.<String>} userClaims
     */
    setUserClaims(userClaims) {
        this['userClaims'] = userClaims;
    }

}

/**
 * @member {String} clientID
 */
V1StaticOIDCIDP.prototype['clientID'] = undefined;

/**
 * @member {Array.<String>} clientScopes
 */
V1StaticOIDCIDP.prototype['clientScopes'] = undefined;

/**
 * @member {String} clientSecret
 */
V1StaticOIDCIDP.prototype['clientSecret'] = undefined;

/**
 * @member {String} issuer
 */
V1StaticOIDCIDP.prototype['issuer'] = undefined;

/**
 * @member {Array.<String>} userClaims
 */
V1StaticOIDCIDP.prototype['userClaims'] = undefined;






export default V1StaticOIDCIDP;

