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
import V1beta1CostElement from './V1beta1CostElement';

/**
 * The V1beta1Cost model module.
 * @module model/V1beta1Cost
 * @version 0.0.1
 */
class V1beta1Cost {
    /**
     * Constructs a new <code>V1beta1Cost</code>.
     * @alias module:model/V1beta1Cost
     */
    constructor() { 
        
        V1beta1Cost.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1beta1Cost</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1beta1Cost} obj Optional instance to populate.
     * @return {module:model/V1beta1Cost} The populated <code>V1beta1Cost</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1beta1Cost();

            if (data.hasOwnProperty('cost')) {
                obj['cost'] = ApiClient.convertToType(data['cost'], 'Number');
            }
            if (data.hasOwnProperty('costElements')) {
                obj['costElements'] = ApiClient.convertToType(data['costElements'], [V1beta1CostElement]);
            }
            if (data.hasOwnProperty('from')) {
                obj['from'] = ApiClient.convertToType(data['from'], 'String');
            }
            if (data.hasOwnProperty('resource')) {
                obj['resource'] = V1Ownership.constructFromObject(data['resource']);
            }
            if (data.hasOwnProperty('resourceIdentifier')) {
                obj['resourceIdentifier'] = ApiClient.convertToType(data['resourceIdentifier'], 'String');
            }
            if (data.hasOwnProperty('retrievedAt')) {
                obj['retrievedAt'] = ApiClient.convertToType(data['retrievedAt'], 'String');
            }
            if (data.hasOwnProperty('team')) {
                obj['team'] = ApiClient.convertToType(data['team'], 'String');
            }
            if (data.hasOwnProperty('teamIdentifier')) {
                obj['teamIdentifier'] = ApiClient.convertToType(data['teamIdentifier'], 'String');
            }
            if (data.hasOwnProperty('to')) {
                obj['to'] = ApiClient.convertToType(data['to'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {Number}
     */
    getCost() {
        return this.cost;
    }

    /**
     * @param {Number} cost
     */
    setCost(cost) {
        this['cost'] = cost;
    }
/**
     * @return {Array.<module:model/V1beta1CostElement>}
     */
    getCostElements() {
        return this.costElements;
    }

    /**
     * @param {Array.<module:model/V1beta1CostElement>} costElements
     */
    setCostElements(costElements) {
        this['costElements'] = costElements;
    }
/**
     * @return {String}
     */
    getFrom() {
        return this.from;
    }

    /**
     * @param {String} from
     */
    setFrom(from) {
        this['from'] = from;
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
    getResourceIdentifier() {
        return this.resourceIdentifier;
    }

    /**
     * @param {String} resourceIdentifier
     */
    setResourceIdentifier(resourceIdentifier) {
        this['resourceIdentifier'] = resourceIdentifier;
    }
/**
     * @return {String}
     */
    getRetrievedAt() {
        return this.retrievedAt;
    }

    /**
     * @param {String} retrievedAt
     */
    setRetrievedAt(retrievedAt) {
        this['retrievedAt'] = retrievedAt;
    }
/**
     * @return {String}
     */
    getTeam() {
        return this.team;
    }

    /**
     * @param {String} team
     */
    setTeam(team) {
        this['team'] = team;
    }
/**
     * @return {String}
     */
    getTeamIdentifier() {
        return this.teamIdentifier;
    }

    /**
     * @param {String} teamIdentifier
     */
    setTeamIdentifier(teamIdentifier) {
        this['teamIdentifier'] = teamIdentifier;
    }
/**
     * @return {String}
     */
    getTo() {
        return this.to;
    }

    /**
     * @param {String} to
     */
    setTo(to) {
        this['to'] = to;
    }

}

/**
 * @member {Number} cost
 */
V1beta1Cost.prototype['cost'] = undefined;

/**
 * @member {Array.<module:model/V1beta1CostElement>} costElements
 */
V1beta1Cost.prototype['costElements'] = undefined;

/**
 * @member {String} from
 */
V1beta1Cost.prototype['from'] = undefined;

/**
 * @member {module:model/V1Ownership} resource
 */
V1beta1Cost.prototype['resource'] = undefined;

/**
 * @member {String} resourceIdentifier
 */
V1beta1Cost.prototype['resourceIdentifier'] = undefined;

/**
 * @member {String} retrievedAt
 */
V1beta1Cost.prototype['retrievedAt'] = undefined;

/**
 * @member {String} team
 */
V1beta1Cost.prototype['team'] = undefined;

/**
 * @member {String} teamIdentifier
 */
V1beta1Cost.prototype['teamIdentifier'] = undefined;

/**
 * @member {String} to
 */
V1beta1Cost.prototype['to'] = undefined;






export default V1beta1Cost;

